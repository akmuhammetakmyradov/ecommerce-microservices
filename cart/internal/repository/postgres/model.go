package postgres

import (
	"cart/internal/models"
)

type DbCartItem struct {
	UserID int64  `db:"user_id"`
	SKU    uint32 `db:"sku"`
	Count  uint32 `db:"count"`
}

func (d DbCartItem) ToDomain() models.CartItem {
	return models.CartItem{
		UserID: d.UserID,
		SKU:    d.SKU,
		Count:  d.Count,
	}
}
