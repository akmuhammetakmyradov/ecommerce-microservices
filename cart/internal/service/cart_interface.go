package service

import (
	"cart/internal/models"
	"context"
)

type CartService interface {
	AddItemToCart(ctx context.Context, params models.CartItem) error
	ListCartItems(ctx context.Context, userID int64) (models.CartItemsList, error)
	DeleteItemFromCart(ctx context.Context, params models.DeleteCartItem) error
	ClearCart(ctx context.Context, userID int64) error
}
