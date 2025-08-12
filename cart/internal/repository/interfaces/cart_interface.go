package interfaces

import (
	"cart/internal/models"
	"context"
)

type CartRepository interface {
	AddItem(ctx context.Context, item models.CartItem) (int64, error)
	CartItemCount(ctx context.Context, userID int64, sku uint32) (uint32, error)
	DeleteCartItem(ctx context.Context, userID int64, sku uint32) error
	ListItems(ctx context.Context, userID int64) ([]models.CartItem, error)
	ClearCart(ctx context.Context, userID int64) error
}
