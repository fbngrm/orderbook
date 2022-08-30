package parse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
)

const (
	BUY  = "buy"
	SELL = "sell"
)

type JSONStreamParser struct {
	reader   io.ReadCloser
	decoder  *json.Decoder
	UpdateCh chan Update
	ErrCh    chan error
}

func NewJSONStreamParser(rc io.ReadCloser) *JSONStreamParser {
	return &JSONStreamParser{
		reader:   rc,
		decoder:  json.NewDecoder(rc),
		UpdateCh: make(chan Update, 1000),
		ErrCh:    make(chan error),
	}
}

func (p *JSONStreamParser) Run(ctx context.Context) (chan Update, chan error) {
	// main loop to read array elements
	go func() {
		hasNext := make(chan bool, 1)
		hasNext <- p.decoder.More()
		for {
			select {
			case <-hasNext:
				err := p.parseNextObject()
				if err != nil {
					p.ErrCh <- err
					return
				}
			case <-ctx.Done():
				p.ErrCh <- ctx.Err()
				return
			default:
				hasNext <- p.decoder.More()
			}
		}
	}()

	return p.UpdateCh, p.ErrCh
}

// Returns and error when EOF is read or the io.Reader is closed.
// Todo, return error when the io.Reader gets closed.
func (p *JSONStreamParser) parseNextObject() error {
	// note, within the loop, we just skip the current iteration/token in case of error.
	// this is not a robust approach and just a proof of concept for the coding challenge.
	token, err := p.decoder.Token()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return err
		}
		fmt.Printf("error reading opening brace for array element: %v\n", err)
		return nil
	}

	// here, we expect either a brace to open an object or a bracket to close the array
	delim, ok := token.(json.Delim)
	if !ok {
		fmt.Printf("expected token to be delimeter: %v\n", token)
		return nil
	}
	// opening the array, we are start parsing
	if delim == '[' {
		return nil
	}
	// closing the array, we are done parsing
	if delim == ']' {
		return nil
	}
	// opening new object
	if delim != '{' {
		fmt.Printf("expected token to be initial opening brace but got: %v\n", delim)
		return nil
	}

	// parse all fields of the object
	for p.decoder.More() {
		p.parseObject()
	}

	closingBrace, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding closing brace: %v\n", err)
		return nil
	}
	if delim, ok := closingBrace.(json.Delim); !ok || delim != '}' {
		fmt.Printf("expected token to be closing brace but got: %v\n", closingBrace)
		return nil
	}
	return nil
}

func (p *JSONStreamParser) parseObject() {
	// parse type key
	key, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding token for type key: %v\n", err)
		return
	}
	// type key
	typeKey, ok := key.(string)
	if !ok {
		fmt.Printf("type key is not type string but: %T\n", key)
		return
	}
	if typeKey != "type" {
		return
	}
	// type value
	val, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding token for type value: %v\n", err)
		return
	}
	typeVal, ok := val.(string)
	if !ok {
		fmt.Println("expect type value to be string")
		return
	}
	if typeVal == "snapshot" {
		p.parseSnapshop()
	}
	if typeVal == "l2update" {
		p.parseUpdate()
	}
}

func (p *JSONStreamParser) parseUpdate() {
	for p.decoder.More() {
		// parse changes
		k, err := p.decoder.Token()
		if err != nil {
			fmt.Printf("error decoding changes token: %v\n", err)
		}
		// type key
		key, ok := k.(string)
		if !ok {
			continue
		}
		if key != "changes" {
			fmt.Printf("expect token value to be changes but got: %q\n", key)
			continue
		}
		p.parseChanges()
	}
}

func (p *JSONStreamParser) parseChanges() {
	// array element
	openingBracket, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding opening bracket for changes array: %v\n", err)
	}
	if delim, ok := openingBracket.(json.Delim); !ok || delim != '[' {
		fmt.Printf("expected opening bracket for changes array but got: %v\n", delim)
	}

	for p.decoder.More() {
		p.parseChange()
	}

	// read closing bracket
	closingBracket, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding closing bracket for changes array: %v\n", err)
	}
	if delim, ok := closingBracket.(json.Delim); !ok || delim != ']' {
		fmt.Printf("expected closing bracket for changes array but got: %v\n", delim)
	}
}

func (p *JSONStreamParser) parseChange() {
	openingBracket, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding opening bracket for change array: %v\n", err)
	}
	if delim, ok := openingBracket.(json.Delim); !ok || delim != '[' {
		fmt.Printf("expected opening bracket for change array but got: %v\n", delim)
	}

	side, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding token for side value: %v\n", err)
	}
	price, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding token for price value: %v\n", err)
	}
	quantity, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding token for quantity value: %v\n", err)
	}

	// note, unchecked type assertions
	p.UpdateCh <- Update{
		Side:     side.(string),
		Price:    price.(string),
		Quantity: quantity.(string),
	}

	closingBracket, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding closing bracket for change array: %v\n", err)
	}
	if delim, ok := closingBracket.(json.Delim); !ok || delim != ']' {
		fmt.Printf("expected closing bracket for change array but got: %v\n", delim)
	}
}

func (p *JSONStreamParser) parseSnapshop() {
	for p.decoder.More() {
		k, err := p.decoder.Token()
		if err != nil {
			fmt.Printf("error decoding token for bids: %v\n", err)
		}
		side, ok := k.(string)
		if !ok {
			continue
		}
		if side == "bids" {
			p.parseBidsOrAsks(BUY)
		}
		if side == "asks" {
			p.parseBidsOrAsks(SELL)
		}

		fmt.Printf("expect token value to be bids or asks but got: %q\n", side)
		continue
	}
}

func (p *JSONStreamParser) parseBidsOrAsks(side string) {
	// array element
	openingBracket, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding opening bracket for bids|asks array: %v\n", err)
	}
	if delim, ok := openingBracket.(json.Delim); !ok || delim != '[' {
		fmt.Printf("expected opening bracket for bids|asks array but got: %v\n", delim)
	}

	for p.decoder.More() {
		p.parseBidOrAskValue(side)
	}

	// read closing bracket
	closingBracket, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding closing bracket for bids|asks array: %v\n", err)
	}
	if delim, ok := closingBracket.(json.Delim); !ok || delim != ']' {
		fmt.Printf("expected closing bracket for bids|asks array but got: %v\n", delim)
	}
}

func (p *JSONStreamParser) parseBidOrAskValue(side string) {
	openingBracket, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding opening bracket for bid|ask array: %v\n", err)
	}
	if delim, ok := openingBracket.(json.Delim); !ok || delim != '[' {
		fmt.Printf("expected opening bracket for bid|ask array but got: %v\n", delim)
	}

	price, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding token for price value: %v\n", err)
	}
	quantity, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding token for quantity value: %v\n", err)
	}

	// note, unchecked type assertions
	p.UpdateCh <- Update{
		Side:     side,
		Price:    price.(string),
		Quantity: quantity.(string),
	}

	closingBracket, err := p.decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	if delim, ok := closingBracket.(json.Delim); !ok || delim != ']' {
		log.Fatal("Expected array closing")
	}
}

// in a future version, we might want to return the ID of the last update we successfully read here.
func (p *JSONStreamParser) Close() error {
	close(p.UpdateCh)
	close(p.ErrCh)
	return p.reader.Close()
}
