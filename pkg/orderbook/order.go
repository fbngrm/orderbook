package orderbook

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Order struct {
	id       string
	side     Side
	quantity decimal.Decimal
	price    decimal.Decimal
}

func NewOrder(orderID string, side Side, quantity, price decimal.Decimal) *Order {
	return &Order{
		id:       orderID,
		side:     side,
		quantity: quantity,
		price:    price,
	}
}

func (o *Order) ID() string {
	return o.id
}

func (o *Order) Side() Side {
	return o.side
}

func (o *Order) Quantity() decimal.Decimal {
	return o.quantity
}

func (o *Order) Price() decimal.Decimal {
	return o.price
}

func (o *Order) String() string {
	return fmt.Sprintf("ID: %s\nSIDE: %s\nQUANTITY: %s\nPRICE: %s\n", o.id, o.side, o.quantity, o.price)
}
