package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/fbngrm/crypto-compare/pkg/parse"
)

func main() {
	input, err := os.Open("./order-book-data.json")
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

	parser := parse.NewJSONStreamParser(input)
	err = parser.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
