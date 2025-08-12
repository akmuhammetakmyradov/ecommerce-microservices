package models

type StockItem struct {
	ID       int64
	UserID   int64
	SKU      uint32
	Name     string
	Type     string
	Count    uint32
	Price    uint32
	Location string
}

type SKU struct {
	SKUID  uint32
	Name   string
	Type   string
	UserID *int64
}

type ListStockParams struct {
	UserID      int64
	Location    string
	PageSize    int64
	CurrentPage int64
}

type ListStock struct {
	Items      []StockItem
	TotalCount int64
	PageNumber int64
	TotalPages int64
}
