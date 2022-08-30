package orderbook

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type input struct {
	orderID  string
	side     Side
	quantity decimal.Decimal
	price    decimal.Decimal
}

type testcase struct {
	name     string
	inputs   []input
	expected map[string]string
	err      error
}

func TestGetSpread(t *testing.T) {

	tests := []testcase{
		{
			name: "",
			inputs: []input{
				{
					orderID:  "01",
					side:     SELL,
					quantity: decimal.NewFromFloat(1.0),
					price:    decimal.NewFromFloat(100.0),
				},
				{
					orderID:  "02",
					side:     SELL,
					quantity: decimal.NewFromFloat(0.1),
					price:    decimal.NewFromFloat(100.1),
				},
				{
					orderID:  "03",
					side:     SELL,
					quantity: decimal.NewFromFloat(4.),
					price:    decimal.NewFromFloat(101.),
				},
				{
					orderID:  "04",
					side:     SELL,
					quantity: decimal.NewFromFloat(7.1),
					price:    decimal.NewFromFloat(102.2),
				},
				{
					orderID:  "05",
					side:     SELL,
					quantity: decimal.NewFromFloat(4.5),
					price:    decimal.NewFromFloat(106.6),
				},
				{
					orderID:  "06",
					side:     SELL,
					quantity: decimal.NewFromFloat(5.8),
					price:    decimal.NewFromFloat(106.7),
				},
				{
					orderID:  "07",
					side:     SELL,
					quantity: decimal.NewFromFloat(8.0),
					price:    decimal.NewFromFloat(107.1),
				},
				{
					orderID:  "07",
					side:     SELL,
					quantity: decimal.NewFromFloat(3.0),
					price:    decimal.NewFromFloat(107.2),
				},
				{
					orderID:  "08",
					side:     SELL,
					quantity: decimal.NewFromFloat(4.0),
					price:    decimal.NewFromFloat(107.3),
				},
				{
					orderID:  "08",
					side:     SELL,
					quantity: decimal.NewFromFloat(1.3),
					price:    decimal.NewFromFloat(108.0),
				},
				{
					orderID:  "09",
					side:     BUY,
					quantity: decimal.NewFromFloat(5.1),
					price:    decimal.NewFromFloat(99.6),
				},
				{
					orderID:  "10",
					side:     BUY,
					quantity: decimal.NewFromFloat(3.4),
					price:    decimal.NewFromFloat(98.7),
				},
				{
					orderID:  "11",
					side:     BUY,
					quantity: decimal.NewFromFloat(4.3),
					price:    decimal.NewFromFloat(99.5),
				},
				{
					orderID:  "12",
					side:     BUY,
					quantity: decimal.NewFromFloat(5.7),
					price:    decimal.NewFromFloat(99.4),
				},
				{
					orderID:  "13",
					side:     BUY,
					quantity: decimal.NewFromFloat(7.1),
					price:    decimal.NewFromFloat(93.4),
				},
				{
					orderID:  "14",
					side:     BUY,
					quantity: decimal.NewFromFloat(21.1),
					price:    decimal.NewFromFloat(96.5),
				},
				{
					orderID:  "15",
					side:     BUY,
					quantity: decimal.NewFromFloat(41.2),
					price:    decimal.NewFromFloat(94.3),
				},
				{
					orderID:  "16",
					side:     BUY,
					quantity: decimal.NewFromFloat(1.4),
					price:    decimal.NewFromFloat(93.1),
				},
				{
					orderID:  "17",
					side:     BUY,
					quantity: decimal.NewFromFloat(1.2),
					price:    decimal.NewFromFloat(96.3),
				},
				{
					orderID:  "18",
					side:     BUY,
					quantity: decimal.NewFromFloat(1.8),
					price:    decimal.NewFromFloat(91.0),
				},
				// update 1
				{
					orderID:  "01",
					side:     SELL,
					quantity: decimal.NewFromFloat(0),
					price:    decimal.NewFromFloat(100.0),
				},
				{
					orderID:  "20",
					side:     SELL,
					quantity: decimal.NewFromFloat(23.1),
					price:    decimal.NewFromFloat(120.0),
				},
				{
					orderID:  "21",
					side:     BUY,
					quantity: decimal.NewFromFloat(12.2),
					price:    decimal.NewFromFloat(80.),
				},
				{
					orderID:  "22",
					side:     BUY,
					quantity: decimal.NewFromFloat(1.4),
					price:    decimal.NewFromFloat(98.7),
				},
				// update 2
				{
					orderID:  "23",
					side:     SELL,
					quantity: decimal.NewFromFloat(6.0),
					price:    decimal.NewFromFloat(101.0),
				},
				{
					orderID:  "24",
					side:     SELL,
					quantity: decimal.NewFromFloat(0.2),
					price:    decimal.NewFromFloat(100.3),
				},
				{
					orderID:  "09",
					side:     BUY,
					quantity: decimal.NewFromFloat(4.4),
					price:    decimal.NewFromFloat(99.6),
				},
			},
			expected: map[string]string{
				`{{"99.6", "5.1"}, {"100.0", "1.0"}}`: `{{"99.6", "5.1"}, {"100.0", "1.0"}}`,
				`{{"99.6", "5.1"}, {"100.1", "0.1"}}`: `{{"99.6", "5.1"}, {"100.1", "0.1"}}`,
				`{{"99.6", "4.4"}, {"100.1", "0.1"}}`: `{{"99.6", "4.4"}, {"100.1", "0.1"}}`,
			},
		},
	}

	ob := NewOrderBook()
	res := make(map[string]string)
	for _, tc := range tests {
		for _, input := range tc.inputs {
			if input.quantity.String() == "0" {
				ob.CancelOrder(input.orderID)
				continue
			}
			err := ob.UpdateOrder(input.orderID, input.side, input.quantity, input.price)
			if err != tc.err {
				t.Logf("unexpected error: %v", err)
			}
			spread, err := ob.GetSpread().MarshalJSON()
			if err != nil {
				t.Log(err)
				t.FailNow()
			}
			spreadStr := string(spread)
			res[spreadStr] = spreadStr
		}

		for key, val := range tc.expected {
			assert.Equal(t, val, res[key])
		}
	}
}
