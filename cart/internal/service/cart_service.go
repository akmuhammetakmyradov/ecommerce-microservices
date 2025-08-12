package service

import (
	"cart/internal/constants"
	"cart/internal/models"
	"cart/internal/repository/interfaces"
	"cart/pkg/log"
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
)

type Service struct {
	repo      interfaces.CartRepository
	stock     interfaces.StockService
	kafkaProd interfaces.KafkaProd
	logger    log.Logger
}

func NewService(repo interfaces.CartRepository, stock interfaces.StockService, kafkaProd interfaces.KafkaProd, logger log.Logger) *Service {
	return &Service{
		repo:      repo,
		stock:     stock,
		kafkaProd: kafkaProd,
		logger:    logger,
	}
}

func (s *Service) AddItemToCart(ctx context.Context, params models.CartItem) error {
	var (
		addedType      = "cart_item_added"
		status         = "success"
		reason         string
		isInsufficient bool
		cartId         int64
	)

	ctx, span := otel.Tracer("cart-service").Start(ctx, "CartService.AddItemToCart")
	defer span.End()

	skuItem, err := s.stock.GetSKU(ctx, params.SKU)
	if err != nil {
		s.logger.Errorf("err in get sku in AddItemToCart: %v", err)

		if errors.Is(err, constants.ErrNotFound) {
			return constants.ErrInvalidSKU
		}

		return fmt.Errorf("failed to validate SKU: %w", err)
	}

	cartItemCount, err := s.repo.CartItemCount(ctx, params.UserID, params.SKU)
	if err != nil {
		s.logger.Errorf("err in CartItemCount: %v", err)
		return err
	}

	if skuItem.Count < params.Count+cartItemCount {
		reason = constants.ErrInsufficientStocks.Error()
		addedType = "cart_item_failed"
		isInsufficient = true
	}

	if !isInsufficient {
		cartId, err = s.repo.AddItem(ctx, params)
		if err != nil {
			s.logger.Errorf("err in AddItem: %v", err)
			return err
		}
	}

	msg, timestamp, err := BuildKafkaEvent(addedType, cartId, skuItem.Price, reason, status, params)
	if err != nil {
		s.logger.Errorf("err in build kafka event: %v", err)
	}

	err = s.kafkaProd.Produce(msg, fmt.Sprint(params.SKU), timestamp)
	if err != nil {
		s.logger.Errorf("err in produce kafka msg: %v", err)
	}

	if isInsufficient {
		return constants.ErrInsufficientStocks
	}

	return nil
}

func (s *Service) ListCartItems(ctx context.Context, userID int64) (models.CartItemsList, error) {
	ctx, span := otel.Tracer("cart-service").Start(ctx, "CartService.ListCartItems")
	defer span.End()

	var result models.CartItemsList
	var total uint32

	items, err := s.repo.ListItems(ctx, userID)
	if err != nil {
		return result, err
	}

	for _, item := range items {
		stockItem, err := s.stock.GetSKU(ctx, item.SKU)
		if err != nil {
			s.logger.Errorf("failed to fetch stock info for SKU %d: %v", item.SKU, err)
			continue
		}

		result.Items = append(result.Items, models.CartItemModel{
			SKU:   item.SKU,
			Count: item.Count,
			Name:  stockItem.Name,
			Price: stockItem.Price,
		})

		total += stockItem.Price * uint32(item.Count)
	}

	result.TotalPrice = total

	return result, nil
}

func (s *Service) DeleteItemFromCart(ctx context.Context, params models.DeleteCartItem) error {
	ctx, span := otel.Tracer("cart-service").Start(ctx, "CartService.DeleteItemFromCart")
	defer span.End()

	err := s.repo.DeleteCartItem(ctx, params.UserID, params.SKU)
	if err != nil && errors.Is(err, constants.ErrNotRowAffected) {
		return constants.ErrNotFound
	}

	return err
}

func (s *Service) ClearCart(ctx context.Context, userID int64) error {
	ctx, span := otel.Tracer("cart-service").Start(ctx, "CartService.ClearCart")
	defer span.End()

	err := s.repo.ClearCart(ctx, userID)
	if err != nil && errors.Is(err, constants.ErrNotRowAffected) {
		return constants.ErrNotFound
	}

	return err
}
