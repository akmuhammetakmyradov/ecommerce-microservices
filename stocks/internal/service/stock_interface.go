package service

import (
	"context"
	"stocks/internal/models"
)

type StockService interface {
	AddItem(ctx context.Context, item models.StockItem) error
	DeleteItem(ctx context.Context, sku uint32) error
	ListByLocation(ctx context.Context, params models.ListStockParams) (models.ListStock, error)
	GetItemBySKU(ctx context.Context, sku uint32) (models.StockItem, error)
}
