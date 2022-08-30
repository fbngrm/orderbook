package orderbook

import (
	rbtx "github.com/emirpasic/gods/examples/redblacktreeextended"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/shopspring/decimal"
)

type OrderSide struct {
	prices    map[string]*Order
	priceTree *rbtx.RedBlackTreeExtended
	numOrders int
	depth     int
}

func comparator(a, b interface{}) int {
	return a.(decimal.Decimal).Cmp(b.(decimal.Decimal))
}

func NewOrderSide() *OrderSide {
	return &OrderSide{
		priceTree: &rbtx.RedBlackTreeExtended{
			Tree: rbt.NewWith(comparator),
		},
		prices: map[string]*Order{},
	}
}

func (os *OrderSide) Append(o *Order) *Order {
	price := o.Price()
	strPrice := price.String()

	_, ok := os.prices[strPrice]
	if !ok {
		os.prices[strPrice] = o
		os.priceTree.Put(price, o)
		os.depth++
	}
	os.numOrders++
	return o
}

func (os *OrderSide) Remove(o *Order) *Order {
	price := o.Price()
	strPrice := price.String()

	if _, ok := os.prices[strPrice]; ok {
		delete(os.prices, strPrice)
		os.priceTree.Remove(price)
		os.depth--
	}
	os.numOrders--
	return o
}

func (os *OrderSide) MaxPriceOrder() *Order {
	if os.depth > 0 {
		if value, found := os.priceTree.GetMax(); found {
			return value.(*Order)
		}
	}
	return nil
}

func (os *OrderSide) MinPriceOrder() *Order {
	if os.depth > 0 {
		if value, found := os.priceTree.GetMin(); found {
			return value.(*Order)
		}
	}
	return nil
}
