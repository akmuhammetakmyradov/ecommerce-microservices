package grpcserver

import (
	cartapi "cart/pkg/api/cart"
	"cart/pkg/log"
	"cart/pkg/metrics"
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

	err := cartapi.RegisterCartServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts)
	if err != nil {
		return nil, err
	}

	metricsWrapped := MetricsMiddleware(mux, m)
	otelHandler := otelhttp.NewHandler(
		metricsWrapped,
		"cart-grpc-gateway",
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		otelhttp.WithSpanOptions(
			trace.WithAttributes(attribute.Bool("suppress-write-header-warning", true)),
		),
	)

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
