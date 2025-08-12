package grpcserver

import (
	"cart/internal/constants"
	cartapi "cart/pkg/api/cart"
)

func ValidateAddItemToCart(req *cartapi.AddItemToCartRequest) error {
	if req.UserId <= 0 {
		return constants.ErrInvalidUserID
	}

	if req.Sku == 0 {
		return constants.ErrInvalidSKU
	}

	if req.Count == 0 {
		return constants.ErrInvalidCount
	}

	return nil
}

func ValidateDeleteItemFromCart(req *cartapi.DeleteItemFromCartRequest) error {
	if req.UserId <= 0 {
		return constants.ErrInvalidUserID
	}

	if req.Sku == 0 {
		return constants.ErrInvalidSKU
	}

	return nil
}

func ValidateCartList(req *cartapi.CartListRequest) error {
	if req.UserId <= 0 {
		return constants.ErrInvalidUserID
	}

	return nil
}

func ValidateClearCart(req *cartapi.ClearCartRequest) error {
	if req.UserId <= 0 {
		return constants.ErrInvalidUserID
	}

	return nil
}
