package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/fbngrm/crypto-compare/pkg/orderbook"
	"github.com/fbngrm/crypto-compare/pkg/parse"
	"github.com/shopspring/decimal"
)

func main() {
	input, err := os.Open("./testdata/order-book-data.json")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	quitCh := make(chan os.Signal, 1)
	// interrupt signal sent from terminal
	signal.Notify(quitCh, os.Interrupt)
	go func() {
		sig := <-quitCh
		log.Printf("received %v, shutting down\n", sig)
		cancel()
	}()

	book := orderbook.NewOrderBook()

	parser := parse.NewJSONStreamParser(input)

	// todo: factor out of main
	updateCh, errCh := parser.Run(ctx)
	go func() {
		for update := range updateCh {
			id := update.Side + update.Price // todo: use a unique hash here for the ID
			price, err := decimal.NewFromString(update.Price)
			if err != nil {
				log.Println(err)
				continue
			}
			quantity, err := decimal.NewFromString(update.Quantity)
			if err != nil {
				log.Println(err)
				continue
			}
			// delete zero orders
			if quantity.String() == "0" {
				book.CancelOrder(id)
				continue
			}
			side, err := orderbook.NewSide(update.Side)
			if err != nil {
				log.Println(err)
				continue
			}
			err = book.UpdateOrder(id, side, quantity, price)
			if err != nil {
				log.Println(err)
				continue
			}
			b, err := book.GetSpread().MarshalJSON()
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(string(b))
		}
	}()

	err = <-errCh // we block until the stream is closed or the context is canceled
	if err != nil {
		// here we could handle EOF or closed input streams and graceful closing of parser in case of error
		fmt.Println(err)
	}

	<-ctx.Done()
	parser.Close()
}
