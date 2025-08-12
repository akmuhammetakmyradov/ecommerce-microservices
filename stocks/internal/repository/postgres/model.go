package postgres

import (
	"stocks/internal/models"
)

type DbStockItem struct {
	ID       int64  `db:"id"`
	UserID   int64  `db:"user_id"`
	SKU      uint32 `db:"sku"`
	Name     string `db:"name"`
	Type     string `db:"type"`
	Count    uint32 `db:"count"`
	Price    uint32 `db:"price"`
	Location string `db:"location"`
}

type DbSKU struct {
	SKUID  uint32 `db:"sku_id"`
	Name   string `db:"name"`
	Type   string `db:"type"`
	UserID *int64 `db:"user_id"`
}

func (d DbStockItem) ToDomain() models.StockItem {
	return models.StockItem{
		ID:       d.ID,
		UserID:   d.UserID,
		SKU:      d.SKU,
		Name:     d.Name,
		Type:     d.Type,
		Count:    d.Count,
		Price:    d.Price,
		Location: d.Location,
	}
}

func (d DbSKU) ToDomain() models.SKU {
	return models.SKU{
		SKUID:  d.SKUID,
		Name:   d.Name,
		Type:   d.Type,
		UserID: d.UserID,
	}
}
