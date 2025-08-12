package service

import (
	"context"
	"errors"
	"fmt"
	"stocks/internal/constants"
	"stocks/internal/models"
	"stocks/internal/repository/interfaces"
	"stocks/pkg/log"

	trm "github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
)

type Service struct {
	repo      interfaces.StockRepository
	tm        trm.Manager
	kafkaProd interfaces.KafkaProd
	logger    log.Logger
}

func NewService(repo interfaces.StockRepository, tm trm.Manager, kafkaProd interfaces.KafkaProd, logger log.Logger) *Service {
	return &Service{
		repo:      repo,
		tm:        tm,
		kafkaProd: kafkaProd,
		logger:    logger,
	}
}

func (s *Service) AddItem(ctx context.Context, item models.StockItem) error {
	ctx, span := otel.Tracer("stock-service").Start(ctx, "StockService.AddItem")
	defer span.End()

	var addedType string

	err := s.tm.Do(ctx, func(ctx context.Context) error {
		sku, err := s.repo.GetSKUByID(ctx, item.SKU)
		if err != nil {
			s.logger.Errorf("err in get sku in AddItem: %v", err)

			if errors.Is(err, pgx.ErrNoRows) {
				return constants.ErrInvalidSKU
			}

			return err
		}

		if sku.UserID != nil && *sku.UserID != item.UserID {
			return constants.ErrAlreadyAdded
		}

		addedType, err = s.repo.AddItem(ctx, item)
		if err != nil {
			s.logger.Errorf("err in add item: %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Errorf("err transaction manager AddItem: %v", err)
		return err
	}

	msg, timestamp, err := BuildKafkaEvent(addedType, item)
	if err != nil {
		s.logger.Errorf("err in build kafka event: %v", err)
	}

	err = s.kafkaProd.Produce(msg, fmt.Sprint(item.SKU), timestamp)
	if err != nil {
		s.logger.Errorf("err in produce kafka msg: %v", err)
	}

	return nil
}

func (s *Service) DeleteItem(ctx context.Context, sku uint32) error {
	ctx, span := otel.Tracer("stock-service").Start(ctx, "StockService.DeleteItem")
	defer span.End()

	err := s.repo.DeleteItem(ctx, sku)
	if err != nil && errors.Is(err, constants.ErrNotRowAffected) {
		return constants.ErrNotFound
	}

	return err
}

func (s *Service) ListByLocation(ctx context.Context, params models.ListStockParams) (models.ListStock, error) {
	ctx, span := otel.Tracer("stock-service").Start(ctx, "StockService.ListByLocation")
	defer span.End()

	var result models.ListStock
	limit := params.PageSize
	offset := (params.CurrentPage - 1) * params.PageSize

	err := s.tm.Do(ctx, func(ctx context.Context) error {
		items, err := s.repo.GetItemsByLocation(ctx, params.Location, params.UserID, limit, offset)
		if err != nil {
			return err
		}

		count, err := s.repo.CountItemsByLocation(ctx, params.Location, params.UserID)
		if err != nil {
			return err
		}

		totalPages := (count + limit - 1) / limit
		result.TotalPages = totalPages
		result.Items = items
		result.TotalCount = count
		result.PageNumber = params.CurrentPage

		return nil
	})

	if err != nil {
		return models.ListStock{}, err
	}

	return result, nil
}

func (s *Service) GetItemBySKU(ctx context.Context, sku uint32) (models.StockItem, error) {
	ctx, span := otel.Tracer("stock-service").Start(ctx, "StockService.GetItemBySKU")
	defer span.End()

	item, err := s.repo.GetItemBySKU(ctx, sku)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.StockItem{}, constants.ErrNotFound
		}

		return models.StockItem{}, err
	}

	return item, nil
}
