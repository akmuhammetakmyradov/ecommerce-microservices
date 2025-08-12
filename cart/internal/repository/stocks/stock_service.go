package stocks

import (
	"cart/internal/constants"
	"cart/internal/models"
	"cart/internal/repository/interfaces"
	stocksapi "cart/pkg/api/stocks"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type grpcStockService struct {
	client  stocksapi.StockServiceClient
	conn    *grpc.ClientConn
	timeout time.Duration
}

func NewGRPCStockService(serverAddr string) (interfaces.StockService, error) {
	conn, err := grpc.NewClient(serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to stock service: %w", err)
	}

	client := stocksapi.NewStockServiceClient(conn)

	return &grpcStockService{
		client:  client,
		conn:    conn,
		timeout: constants.ReadTimeout,
	}, nil
}

func (s *grpcStockService) GetSKU(ctx context.Context, sku uint32) (models.StockItem, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	resp, err := s.client.GetStock(ctx, &stocksapi.GetStockRequest{Sku: sku})

	if err != nil {
		st, ok := status.FromError(err)
		switch {
		case ok && st.Code() == codes.NotFound:
			return models.StockItem{}, constants.ErrNotFound
		case ctx.Err() == context.DeadlineExceeded:
			return models.StockItem{}, fmt.Errorf("stock service request timed out")
		default:
			return models.StockItem{}, fmt.Errorf("gRPC stock service error: %w", err)
		}
	}

	if resp.Stock == nil {
		return models.StockItem{}, constants.ErrNotFound
	}

	return models.StockItem{
		SKU:      resp.Stock.Sku,
		Name:     resp.Stock.Name,
		Type:     resp.Stock.Type,
		Count:    resp.Stock.Count,
		Price:    resp.Stock.Price,
		Location: resp.Stock.Location,
	}, nil
}

func (s *grpcStockService) Close() error {
	return s.conn.Close()
}
