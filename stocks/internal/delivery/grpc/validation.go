package grpcserver

import (
	"errors"
	stocksapi "stocks/pkg/api/stocks"
)

func ValidateAddStock(req *stocksapi.AddStockRequest) error {
	if req.UserId == 0 {
		return errors.New("user_id is required")
	}

	if req.Sku == 0 {
		return errors.New("sku is required")
	}

	if req.Count == 0 {
		return errors.New("count must be greater than 0")
	}

	if req.Location == "" {
		return errors.New("location is required")
	}

	return nil
}

func ValidateDeleteStock(req *stocksapi.DeleteStockRequest) error {
	if req.UserId == 0 {
		return errors.New("userId is required")
	}

	if req.Sku == 0 {
		return errors.New("sku is required")
	}

	return nil
}

func ValidateListStocks(req *stocksapi.ListStocksByLocationRequest) error {
	if req.UserId == 0 {
		return errors.New("userId is required")
	}

	if req.Location == "" {
		return errors.New("location is required")
	}

	if req.PageSize <= 0 {
		return errors.New("pageSize must be greater than 0")
	}

	if req.CurrentPage <= 0 {
		return errors.New("currentPage must be greater than 0")
	}

	return nil
}

func ValidateGetStock(req *stocksapi.GetStockRequest) error {
	if req.Sku == 0 {
		return errors.New("SKU must be greater than 0")
	}

	return nil
}
