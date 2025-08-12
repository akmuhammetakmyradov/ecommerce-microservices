package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"stocks/internal/config"
	"stocks/internal/constants"
	"stocks/internal/migrations"
	kconstructor "stocks/internal/repository/kafka"
	"stocks/internal/repository/postgres"
	"stocks/internal/service"
	"stocks/pkg/log"
	"stocks/pkg/log/zap"
	"stocks/pkg/metrics"
	"stocks/pkg/postgresql"
	"stocks/pkg/tracing"
	"syscall"

	grpcserver "stocks/internal/delivery/grpc"

	tmsql "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	trm "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"google.golang.org/grpc"
)

type App struct {
	db             postgresql.Client
	cfg            *config.Configs
	grpcServer     *grpc.Server
	gateway        grpcserver.Gateway
	logger         log.Logger
	shutdownTracer func(context.Context) error
	metricsServer  metrics.MetricsServer
}

func NewApp(ctx context.Context) (*App, error) {
	cfg := config.GetConfig()

	logger, err := zap.NewLogger(cfg.Listen.ServiceName, cfg.Listen.Env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		return nil, err
	}

	shutdownTracer, err := tracing.InitTracer(cfg.Listen.ServiceName, cfg.Tracing.JaegerEndpoint)
	if err != nil {
		logger.Errorf("failed to init tracer: %v", err)
		return nil, err
	}

	stockMetrics, err := metrics.RegisterMetrics()
	if err != nil {
		logger.Errorf("Failed to register metrics: %v", err)
		return nil, err
	}

	metricsServer := metrics.NewServer(stockMetrics, cfg.Metrics.Port, logger)

	db, err := postgresql.NewPostgres(ctx, cfg)
	if err != nil {
		logger.Errorf("❌ Failed to connect to DB: %v", err)
		return nil, err
	}

	err = migrations.RunManualMigrations(ctx, db, cfg.Listen.MigrationsPath)
	if err != nil {
		logger.Errorf("failed to run migrations: %v", err)
		return nil, err
	}

	driver := tmsql.NewDefaultFactory(db)
	tm := trm.Must(driver)

	kafkaProd, err := kconstructor.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	if err != nil {
		logger.Errorf("failed to connect kafka: %v", err)
		return nil, err
	}

	repo := postgres.NewRepository(db, tmsql.DefaultCtxGetter)
	svc := service.NewService(repo, tm, kafkaProd, logger)

	// gRPC Server Setup
	grpcServer := grpcserver.NewGRPCServer(svc, logger)

	// gRPC-Gateway Setup
	gateway, err := grpcserver.NewGateway(ctx, cfg.Listen.GRPCPort, cfg.Listen.GatewayPort, logger, stockMetrics)
	if err != nil {
		logger.Errorf("failed to create gateway: %v", err)
		return nil, err
	}

	return &App{
		db:             db,
		cfg:            cfg,
		grpcServer:     grpcServer,
		gateway:        gateway,
		logger:         logger,
		shutdownTracer: shutdownTracer,
		metricsServer:  metricsServer,
	}, nil
}

func (a *App) Run() error {
	defer a.logger.Close()

	serverErrors := make(chan error, 3)

	// Start metrics server
	go func() {
		if err := a.metricsServer.Run(); err != nil {
			serverErrors <- fmt.Errorf("metrics server error: %w", err)
		}
	}()

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":"+a.cfg.Listen.GRPCPort)
		if err != nil {
			serverErrors <- fmt.Errorf("failed to listen gRPC: %w", err)
			return
		}

		a.logger.Infof("✅ gRPC server is running on port %s\n", a.cfg.Listen.GRPCPort)
		if err := a.grpcServer.Serve(lis); err != nil {
			serverErrors <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	// Start Gateway server
	go func() {
		a.logger.Infof("✅ gRPC-Gateway is running on port %s\n", a.cfg.Listen.GatewayPort)
		if err := a.gateway.Run(); err != nil {
			serverErrors <- fmt.Errorf("gateway error: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-quit:
		a.logger.Infof("⚠️ Received signal: %s. Shutting down...\n", sig)
	case err := <-serverErrors:
		return errors.New("server failed to start or stopped unexpectedly: " + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.ServerTimeout)
	defer cancel()

	defer func() {
		if err := a.shutdownTracer(ctx); err != nil {
			a.logger.Errorf("failed to shutdown tracer: %v", err)
		}
	}()

	var shutdownErrors []error

	// Shutdown gRPC server
	a.grpcServer.GracefulStop()
	a.logger.Info("✅ gRPC server shutdown complete")

	// Shutdown Gateway
	if err := a.gateway.Shutdown(ctx); err != nil {
		shutdownErrors = append(shutdownErrors, fmt.Errorf("gateway shutdown failed: %w", err))
	} else {
		a.logger.Info("✅ Gateway shutdown complete")
	}

	// Shutdown metrics server
	if err := a.metricsServer.Shutdown(ctx); err != nil {
		shutdownErrors = append(shutdownErrors, err)
	} else {
		a.logger.Info("✅ Metrics server shutdown complete")
	}

	if a.db != nil {
		a.db.Close()
		a.logger.Info("✅ Database connection closed")
	}

	if len(shutdownErrors) > 0 {
		return fmt.Errorf("shutdown completed with errors: %v", shutdownErrors)
	}

	a.logger.Info("✅ Graceful shutdown complete")

	return nil
}

func (a *App) Logger() log.Logger {
	return a.logger
}
