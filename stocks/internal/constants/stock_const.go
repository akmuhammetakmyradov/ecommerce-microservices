package constants

import (
	"errors"
	"time"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrInvalidSKU     = errors.New("invalid sku")
	ErrNotRowAffected = errors.New("not row affected")
	ErrAlreadyAdded   = errors.New("already added")
	ErrUnknownType    = errors.New("unknown event type")
)

const (
	InternalServerErrMessage = "Something went wrong in server!"
	ServerTimeout            = 5 * time.Second
	ReadTimeout              = 3 * time.Second
)
