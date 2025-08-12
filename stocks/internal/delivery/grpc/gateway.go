package grpcserver

import (
	"context"
	"net/http"
	"stocks/pkg/log"
	"stocks/pkg/metrics"

	stocksapi "stocks/pkg/api/stocks"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	server  *http.Server
	metrics metrics.Metrics
}

type Gateway interface {
	Run() error
	Shutdown(ctx context.Context) error
}

func NewGateway(ctx context.Context, grpcPort, gatewayPort string, logger log.Logger, m metrics.Metrics) (Gateway, error) {
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(ErrorMiddleware(logger, m)),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := stocksapi.RegisterStockServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts)
	if err != nil {
		return nil, err
	}

	metricsWrapped := MetricsMiddleware(mux, m)
	otelHandler := otelhttp.NewHandler(metricsWrapped, "stocks-grpc-gateway")

	return &Server{
		server: &http.Server{
			Addr:    ":" + gatewayPort,
			Handler: otelHandler,
		},
		metrics: m,
	}, nil
}

func (g *Server) Run() error {
	return g.server.ListenAndServe()
}

func (g *Server) Shutdown(ctx context.Context) error {
	return g.server.Shutdown(ctx)
}
