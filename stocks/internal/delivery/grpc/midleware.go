package grpcserver

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"stocks/pkg/log"
	"stocks/pkg/metrics"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.status != 0 {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func MetricsMiddleware(next http.Handler, m metrics.Metrics) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w}

		defer func() {
			duration := time.Since(start).Seconds()
			statusCode := rw.status

			if statusCode == 0 {
				statusCode = http.StatusOK
			}

			m.ObserveLatency(r.URL.Path, r.Method, strconv.Itoa(statusCode), duration)
		}()

		next.ServeHTTP(rw, r)
	})
}

func ErrorMiddleware(logger log.Logger, m metrics.Metrics) runtime.ErrorHandlerFunc {
	return func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
		s, ok := status.FromError(err)
		if !ok {
			s = status.New(codes.Unknown, err.Error())
		}

		httpStatus := runtime.HTTPStatusFromCode(s.Code())
		m.IncError(r.URL.Path, r.Method, strconv.Itoa(httpStatus))

		var traceID string
		if span := trace.SpanFromContext(ctx); span != nil && span.SpanContext().IsValid() {
			traceID = span.SpanContext().TraceID().String()
		}

		logger.Error("gRPC-Gateway error",
			log.String("method", r.Method),
			log.String("path", r.URL.Path),
			log.Int("status", httpStatus),
			log.String("trace_id", traceID),
			log.String("error", s.Message()),
		)

		w.Header().Set("Content-Type", marshaler.ContentType("application/json"))
		w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))

		response := map[string]interface{}{
			"error": s.Message(),
			"code":  s.Code(),
		}

		if details := s.Details(); len(details) > 0 {
			response["details"] = details
		}

		if err := marshaler.NewEncoder(w).Encode(response); err != nil {
			logger.Errorf("Failed to encode error response: %v", err)
		}
	}
}
