package grpcserver

import (
	"stocks/internal/models"
	stocksapi "stocks/pkg/api/stocks"
)

func ToAddStockModel(req *stocksapi.AddStockRequest) models.StockItem {
	return models.StockItem{
		UserID:   req.UserId,
		SKU:      req.Sku,
		Count:    req.Count,
		Price:    req.Price,
		Location: req.Location,
	}
}

func ToListStocksModel(req *stocksapi.ListStocksByLocationRequest) models.ListStockParams {
	return models.ListStockParams{
		UserID:      req.UserId,
		Location:    req.Location,
		PageSize:    req.PageSize,
		CurrentPage: req.CurrentPage,
	}
}

func ToListStocksResponse(domain models.ListStock) *stocksapi.ListStocksByLocationResponse {
	items := make([]*stocksapi.StockItem, 0, len(domain.Items))

	for _, item := range domain.Items {
		items = append(items, &stocksapi.StockItem{
			Sku:      item.SKU,
			Name:     item.Name,
			Type:     item.Type,
			Count:    item.Count,
			Price:    item.Price,
			Location: item.Location,
		})
	}

	return &stocksapi.ListStocksByLocationResponse{
		Items:      items,
		TotalCount: domain.TotalCount,
		PageNumber: domain.PageNumber,
		TotalPages: domain.TotalPages,
	}
}
