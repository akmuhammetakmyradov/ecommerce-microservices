package constants

import (
	"errors"
	"time"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrNotRowAffected     = errors.New("not row affected")
	ErrInvalidSKU         = errors.New("invalid sku")
	ErrInvalidCount       = errors.New("count must be greater than 0")
	ErrInvalidUserID      = errors.New("userID must be greater than 0")
	ErrInsufficientStocks = errors.New("insufficient stocks")
	ErrUnknownType        = errors.New("unknown event type")
)

const (
	InternalServerErrMessage = "Something went wrong in server!"
	ServerTimeout            = 5 * time.Second
	ReadTimeout              = 3 * time.Second
)
