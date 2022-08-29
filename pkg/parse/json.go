package parse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
)

type JSONStreamParser struct {
	reader  io.ReadCloser
	decoder *json.Decoder
}

func NewJSONStreamParser(rc io.ReadCloser) *JSONStreamParser {
	return &JSONStreamParser{
		reader:  rc,
		decoder: json.NewDecoder(rc),
	}
}

func (p *JSONStreamParser) Run(ctx context.Context) error {
	// read open bracket
	openingBracket, err := p.decoder.Token()
	if err != nil {
		return fmt.Errorf("error reading initial opening bracket: %w", err)
	}
	if delim, ok := openingBracket.(json.Delim); !ok || delim != '[' {
		return fmt.Errorf("expected token to be initial opening bracket but got: %v", openingBracket)
	}

	// main loop to read array elements and closing bracket
	hasNext := make(chan bool, 1)
	hasNext <- p.parseNextObject()
	for {
		select {
		case <-hasNext:
			hasNext <- p.parseNextObject()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (p *JSONStreamParser) parseNextObject() bool {
	// note, within the loop, we just skip the current iteration/token in case of error.
	// this is not a robust approach and just a proof of concept for the coding challenge.
	token, err := p.decoder.Token()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return false
		}
		fmt.Printf("error reading opening brace for array element: %v\n", err)
		return p.decoder.More()
	}

	// here, we expect either a brace to open an object or a bracket to close the array
	delim, ok := token.(json.Delim)
	if !ok {
		fmt.Printf("expected token to be delimeter: %v\n", token)
		return p.decoder.More()
	}
	// closing the array, we are done parsing
	if delim == ']' {
		return p.decoder.More()
	}
	// opening new object
	if delim != '{' {
		fmt.Printf("expected token to be initial opening brace but got: %v\n", delim)
		return p.decoder.More()
	}

	// parse all fields of the object
	for p.decoder.More() {
		p.parseObject()
	}

	closingBrace, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding closing brace: %v\n", err)
		return p.decoder.More()
	}
	if delim, ok := closingBrace.(json.Delim); !ok || delim != '}' {
		fmt.Printf("expected token to be closing brace but got: %v\n", closingBrace)
		return p.decoder.More()
	}
	return p.decoder.More()
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
	amount, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding token for amount value: %v\n", err)
	}
	number, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding token for number value: %v\n", err)
	}
	fmt.Println(side, amount, number)

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
		// parse bids
		k, err := p.decoder.Token()
		if err != nil {
			fmt.Printf("error decoding token for bids: %v\n", err)
		}
		// type key
		key, ok := k.(string)
		if !ok {
			continue
		}
		if key != "bids" && key != "asks" {
			fmt.Printf("expect token value to be bids or asks but got: %q\n", key)
			continue
		}
		p.parseBidsOrAsks()
	}
}

func (p *JSONStreamParser) parseBidsOrAsks() {
	// array element
	openingBracket, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding opening bracket for bids|asks array: %v\n", err)
	}
	if delim, ok := openingBracket.(json.Delim); !ok || delim != '[' {
		fmt.Printf("expected opening bracket for bids|asks array but got: %v\n", delim)
	}

	for p.decoder.More() {
		p.parseBidOrAskValue()
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

func (p *JSONStreamParser) parseBidOrAskValue() {
	openingBracket, err := p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding opening bracket for bid|ask array: %v\n", err)
	}
	if delim, ok := openingBracket.(json.Delim); !ok || delim != '[' {
		fmt.Printf("expected opening bracket for bid|ask array but got: %v\n", delim)
	}

	_, err = p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding token for amount value: %v\n", err)
	}
	_, err = p.decoder.Token()
	if err != nil {
		fmt.Printf("error decoding token for volume value: %v\n", err)
	}
	// fmt.Println(k, v)

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
	return p.reader.Close()
}
