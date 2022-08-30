package orderbook

import (
	"fmt"
)

type Spread struct {
	highestBidPrice  string
	highestBidAmount string
	lowestAskPrice   string
	lowestAskAmount  string
}

func (s *Spread) MarshalJSON() ([]byte, error) {
	type Bid struct {
		Price  string `json:"highestBidPriceString"`
		Amount string `json:"highestBidAmountString"`
	}
	type Ask struct {
		Price  string `json:"lowestAskPriceString"`
		Amount string `json:"lowestAskAmountString"`
	}

	b := Bid{
		Price:  s.highestBidPrice,
		Amount: s.highestBidAmount,
	}
	a := Ask{
		Price:  s.lowestAskPrice,
		Amount: s.lowestAskAmount,
	}

	return []byte(fmt.Sprintf(
		`{{"%s", "%s"}, {"%s", "%s"}}`,
		b.Price,
		b.Amount,
		a.Price,
		a.Amount,
	)), nil
}
