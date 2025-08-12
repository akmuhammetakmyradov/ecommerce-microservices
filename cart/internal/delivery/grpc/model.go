package grpcserver

import (
	"cart/internal/models"
	cartapi "cart/pkg/api/cart"
)

func ToAddItemCartModel(req *cartapi.AddItemToCartRequest) models.CartItem {
	return models.CartItem{
		UserID: req.UserId,
		SKU:    req.Sku,
		Count:  req.Count,
	}
}

func ToDeleteCartItemModel(req *cartapi.DeleteItemFromCartRequest) models.DeleteCartItem {
	return models.DeleteCartItem{
		UserID: req.UserId,
		SKU:    req.Sku,
	}
}

func ToCartListResponse(domain models.CartItemsList) *cartapi.CartListResponse {
	items := make([]*cartapi.StockItem, 0, len(domain.Items))

	for _, item := range domain.Items {
		items = append(items, &cartapi.StockItem{
			Sku:   item.SKU,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		})
	}

	return &cartapi.CartListResponse{
		Items:      items,
		TotalPrice: domain.TotalPrice,
	}
}
