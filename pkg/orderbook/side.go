package orderbook

import (
	"fmt"
	"strings"
)

type Side int

// SELL (asks) or BUY (bids)
const (
	SELL Side = iota
	BUY
	INVALID
)

func NewSide(s string) (Side, error) {
	if strings.ToLower(s) == SELL.String() {
		return SELL, nil
	}
	if strings.ToLower(s) == BUY.String() {
		return BUY, nil
	}
	return INVALID, fmt.Errorf("side not supported: %q", s)
}

func (s Side) String() string {
	if s == BUY {
		return "buy"
	}
	return "sell"
}
