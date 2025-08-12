package metrics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"stocks/pkg/log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	metricsServer *http.Server
}

type MetricsServer interface {
	Run() error
	Shutdown(ctx context.Context) error
}

func NewServer(m Metrics, metricsPort int64, logger log.Logger) MetricsServer {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			m.IncError(r.URL.Path, r.Method, strconv.Itoa(http.StatusInternalServerError))
			logger.Errorf("Failed to write health check response: %v", err)
		}
	})

	return &Server{
		metricsServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", metricsPort),
			Handler: mux,
		},
	}
}

func (s *Server) Run() error {
	if err := s.metricsServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("metrics server error: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.metricsServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("metrics server shutdown error: %w", err)
	}

	return nil
}
