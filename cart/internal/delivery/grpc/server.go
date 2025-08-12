package grpcserver

import (
	"cart/internal/constants"
	"cart/internal/service"
	cartapi "cart/pkg/api/cart"
	"cart/pkg/log"
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type grpcServer struct {
	cartapi.UnimplementedCartServiceServer
	service service.CartService
	logger  log.Logger
}

func NewGRPCServer(svc service.CartService, logger log.Logger) *grpc.Server {
	grpcServer := &grpcServer{
		service: svc,
		logger:  logger,
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpcLoggingInterceptor(logger)),
	)

	reflection.Register(srv)

	cartapi.RegisterCartServiceServer(srv, grpcServer)

	return srv
}

func (s *grpcServer) AddItemToCart(ctx context.Context, req *cartapi.AddItemToCartRequest) (*cartapi.AddItemToCartResponse, error) {
	ctx, span := otel.Tracer("cart-handler").Start(ctx, "grpcServer.AddItemToCart")
	defer span.End()

	if err := ValidateAddItemToCart(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.service.AddItemToCart(ctx, ToAddItemCartModel(req))
	if err != nil {
		if errors.Is(err, constants.ErrInsufficientStocks) || errors.Is(err, constants.ErrInvalidSKU) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, constants.InternalServerErrMessage)
	}

	return &cartapi.AddItemToCartResponse{Message: "item succesfully added"}, nil
}

func (s *grpcServer) DeleteItemFromCart(ctx context.Context, req *cartapi.DeleteItemFromCartRequest) (*cartapi.DeleteItemFromCartResponse, error) {
	ctx, span := otel.Tracer("cart-handler").Start(ctx, "grpcServer.DeleteItemFromCart")
	defer span.End()

	if err := ValidateDeleteItemFromCart(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.service.DeleteItemFromCart(ctx, ToDeleteCartItemModel(req))
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, constants.InternalServerErrMessage)
	}

	return &cartapi.DeleteItemFromCartResponse{Message: "Stock deleted successfully"}, nil
}

func (s *grpcServer) CartList(ctx context.Context, req *cartapi.CartListRequest) (*cartapi.CartListResponse, error) {
	ctx, span := otel.Tracer("cart-handler").Start(ctx, "grpcServer.CartList")
	defer span.End()

	if err := ValidateCartList(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := s.service.ListCartItems(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, constants.InternalServerErrMessage)
	}

	return ToCartListResponse(result), nil
}

func (s *grpcServer) ClearCart(ctx context.Context, req *cartapi.ClearCartRequest) (*cartapi.ClearCartResponse, error) {
	ctx, span := otel.Tracer("cart-handler").Start(ctx, "grpcServer.ClearCart")
	defer span.End()

	if err := ValidateClearCart(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.service.ClearCart(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, constants.InternalServerErrMessage)
	}

	return &cartapi.ClearCartResponse{Message: "cart succesfully cleared"}, nil
}
