package orderbook

import (
	"encoding/json"
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

	bid, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	ask, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("{%s,%s}", bid, ask)), nil
}
