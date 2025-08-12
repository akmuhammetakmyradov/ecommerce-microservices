package interfaces

import (
	"cart/internal/models"
	"context"
)

type StockService interface {
	GetSKU(ctx context.Context, sku uint32) (models.StockItem, error)
	Close() error
}
