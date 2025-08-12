package grpcserver

import (
	"context"
	"errors"
	"stocks/internal/constants"
	"stocks/internal/service"
	stocksapi "stocks/pkg/api/stocks"
	"stocks/pkg/log"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type grpcServer struct {
	stocksapi.UnimplementedStockServiceServer
	service service.StockService
	logger  log.Logger
}

func NewGRPCServer(svc service.StockService, logger log.Logger) *grpc.Server {
	grpcServer := &grpcServer{
		service: svc,
		logger:  logger,
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpcLoggingInterceptor(logger)),
	)

	reflection.Register(srv)

	stocksapi.RegisterStockServiceServer(srv, grpcServer)

	return srv
}

func (s *grpcServer) AddStock(ctx context.Context, req *stocksapi.AddStockRequest) (*stocksapi.AddStockResponse, error) {
	ctx, span := otel.Tracer("stocks-handler").Start(ctx, "grpcServer.AddStock")
	defer span.End()

	if err := ValidateAddStock(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.service.AddItem(ctx, ToAddStockModel(req))
	if err != nil {
		if errors.Is(err, constants.ErrAlreadyAdded) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		} else if errors.Is(err, constants.ErrInvalidSKU) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, constants.InternalServerErrMessage)
	}

	return &stocksapi.AddStockResponse{Message: "Stock item succesfully created"}, nil
}

func (s *grpcServer) DeleteStock(ctx context.Context, req *stocksapi.DeleteStockRequest) (*stocksapi.DeleteStockResponse, error) {
	ctx, span := otel.Tracer("stocks-handler").Start(ctx, "grpcServer.DeleteStock")
	defer span.End()

	if err := ValidateDeleteStock(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.service.DeleteItem(ctx, req.Sku)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, constants.InternalServerErrMessage)
	}

	return &stocksapi.DeleteStockResponse{Message: "Stock deleted successfully"}, nil
}

func (s *grpcServer) ListStocksByLocation(ctx context.Context, req *stocksapi.ListStocksByLocationRequest) (*stocksapi.ListStocksByLocationResponse, error) {
	ctx, span := otel.Tracer("stocks-handler").Start(ctx, "grpcServer.ListStocksByLocation")
	defer span.End()

	if err := ValidateListStocks(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := s.service.ListByLocation(ctx, ToListStocksModel(req))
	if err != nil {
		return nil, status.Error(codes.Internal, constants.InternalServerErrMessage)
	}

	return ToListStocksResponse(result), nil
}

func (s *grpcServer) GetStock(ctx context.Context, req *stocksapi.GetStockRequest) (*stocksapi.GetStockResponse, error) {
	ctx, span := otel.Tracer("stocks-handler").Start(ctx, "grpcServer.GetStock")
	defer span.End()

	if err := ValidateGetStock(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	stock, err := s.service.GetItemBySKU(ctx, req.Sku)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, constants.InternalServerErrMessage)
	}

	return &stocksapi.GetStockResponse{
		Stock: &stocksapi.StockItem{
			Sku:      stock.SKU,
			Name:     stock.Name,
			Type:     stock.Type,
			Count:    stock.Count,
			Price:    stock.Price,
			Location: stock.Location,
		},
	}, nil
}
