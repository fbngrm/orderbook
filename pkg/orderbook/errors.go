package orderbook

import "errors"

var (
	ErrInvalid         = errors.New("invalid order")
	ErrInvalidQuantity = errors.New("invalid order quantity")
	ErrInvalidPrice    = errors.New("invalid order price")
	ErrOrderExists     = errors.New("order already exists")
)
