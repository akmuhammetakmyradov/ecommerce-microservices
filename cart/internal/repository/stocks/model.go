package stocks

import "cart/internal/models"

type APIResponse[T any] struct {
	Data T `json:"data"`
}

type StockItemResponse struct {
	SKU      uint32 `json:"sku"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Count    uint32 `json:"count"`
	Price    uint32 `json:"price"`
	Location string `json:"location"`
}

type StockItemRequest struct {
	SKU uint32 `json:"sku"`
}

var ErrorResponse struct {
	Message string `json:"message"`
}

func (res StockItemResponse) ToModel() models.StockItem {
	return models.StockItem{
		SKU:      res.SKU,
		Name:     res.Name,
		Type:     res.Type,
		Count:    res.Count,
		Price:    res.Price,
		Location: res.Location,
	}
}
