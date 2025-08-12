package interfaces

import (
	"context"
	"stocks/internal/models"
)

type StockRepository interface {
	AddItem(ctx context.Context, item models.StockItem) (string, error)
	DeleteItem(ctx context.Context, sku uint32) error
	GetItemsByLocation(ctx context.Context, location string, userID, limit, offset int64) ([]models.StockItem, error)
	CountItemsByLocation(ctx context.Context, location string, userID int64) (int64, error)
	GetItemBySKU(ctx context.Context, sku uint32) (models.StockItem, error)
	GetSKUByID(ctx context.Context, skuID uint32) (models.SKU, error)
}
