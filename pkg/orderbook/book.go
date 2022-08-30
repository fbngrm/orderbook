package orderbook

import (
	"github.com/shopspring/decimal"
)

type OrderBook struct {
	orders map[string]*Order
	asks   *OrderSide
	bids   *OrderSide
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		orders: make(map[string]*Order),
		bids:   NewOrderSide(),
		asks:   NewOrderSide(),
	}
}

// IsInvalid returns true if the lowest ask/sell is less than the highest bid/buy.
func (ob *OrderBook) IsInvalid(o *Order) bool {
	if o.Side() == BUY { // bid
		minAsk := ob.asks.MinPriceOrder()
		if minAsk == nil {
			return false
		}
		if o.Price().GreaterThanOrEqual(minAsk.Price()) { // sell
			return true
		}
		return false
	}

	maxBid := ob.bids.MaxPriceOrder()
	if maxBid == nil {
		return false
	}
	if o.Price().LessThan(maxBid.Price()) {
		return true // sell
	}
	return false
}

func (ob *OrderBook) UpdateOrder(orderID string, side Side, quantity, price decimal.Decimal) error {
	if o, ok := ob.orders[orderID]; ok {
		ob.CancelOrder(o.ID())
	}
	return ob.AddOrder(orderID, side, quantity, price)
}

func (ob *OrderBook) AddOrder(orderID string, side Side, quantity, price decimal.Decimal) error {
	if _, ok := ob.orders[orderID]; ok {
		return ErrOrderExists
	}
	if quantity.Sign() <= 0 {
		return ErrInvalidQuantity
	}
	if price.Sign() <= 0 {
		return ErrInvalidPrice
	}

	o := NewOrder(orderID, side, quantity, price)
	if ob.IsInvalid(o) {
		return ErrInvalid
	}

	if side == BUY {
		ob.orders[orderID] = ob.bids.Append(o)
	} else {
		ob.orders[orderID] = ob.asks.Append(o)
	}
	return nil
}

func (ob *OrderBook) CancelOrder(orderID string) *Order {
	e, ok := ob.orders[orderID]
	if !ok {
		return nil
	}

	delete(ob.orders, orderID)

	if e.Side() == BUY {
		return ob.bids.Remove(e)
	}
	return ob.asks.Remove(e)
}

func (ob *OrderBook) GetSpread() *Spread {
	lowestAskPrice := "0"
	lowestAskQuantity := "0"
	minAsk := ob.asks.MinPriceOrder()
	if minAsk != nil {
		lowestAskPrice = minAsk.Price().StringFixed(1)
		lowestAskQuantity = minAsk.Quantity().StringFixed(1)
	}

	highestBidPrice := "0"
	highestBidQuantity := "0"
	maxBid := ob.bids.MaxPriceOrder()
	if maxBid != nil {
		highestBidPrice = maxBid.Price().StringFixed(1)
		highestBidQuantity = maxBid.Quantity().StringFixed(1)
	}

	return &Spread{
		highestBidPrice:  highestBidPrice,
		highestBidAmount: highestBidQuantity,
		lowestAskPrice:   lowestAskPrice,
		lowestAskAmount:  lowestAskQuantity,
	}
}

func (ob *OrderBook) String() string {
	s := "------------------\n"
	for _, o := range ob.orders {
		s += o.String()
		s += "------------------\n"
	}
	return s
}
