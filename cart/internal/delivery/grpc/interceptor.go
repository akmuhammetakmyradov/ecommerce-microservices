package grpcserver

import (
	"cart/pkg/log"
	"context"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func grpcLoggingInterceptor(logger log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, span := otel.Tracer("cart-grpc-server").Start(ctx, info.FullMethod)
		defer span.End()

		resp, err := handler(ctx, req)

		var traceID string
		if span.SpanContext().IsValid() {
			traceID = span.SpanContext().TraceID().String()
		}

		if err != nil {
			st, ok := status.FromError(err)
			statusCode := codes.Unknown
			errMsg := err.Error()

			if ok {
				statusCode = st.Code()
				errMsg = st.Message()
			}

			logger.Error("gRPC call failed",
				log.String("method", info.FullMethod),
				log.Any("status", statusCode),
				log.String("trace_id", traceID),
				log.String("error", errMsg),
			)
		}

		return resp, err
	}
}
