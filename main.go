package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Snapshot struct {
	Type string     `json:"type"`
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

func main() {
	input, err := os.Open("./order-book-data.json")
	if err != nil {
		log.Fatal(err)
	}
	dec := json.NewDecoder(input)

	// read open bracket
	openingBracket, err := dec.Token()

	if err != nil {
		log.Fatal(err)
	}
	if delim, ok := openingBracket.(json.Delim); !ok || delim != '[' {
		log.Fatal("Expected array")
	}

	// read strings one by one
	for dec.More() {
		// array element
		openingBrace, err := dec.Token()
		if err != nil {
			log.Fatal(err)
		}
		if delim, ok := openingBrace.(json.Delim); !ok || delim != '{' {
			log.Fatal("Expected object")
		}

		typeVal := ""
		for dec.More() {
			// parse type
			k, err := dec.Token()
			if err != nil {
				log.Println(err)
				continue
			}
			// type key
			typeKey, ok := k.(string)
			if !ok {
				continue
			}
			if typeKey != "type" {
				continue
			}
			// type value
			v, err := dec.Token()
			if err != nil {
				log.Println(err)
			}
			typeVal, ok = v.(string)
			if !ok {
				log.Println("expect type value to be string")
				continue
			}
			if typeVal == "snapshot" {
				fmt.Println(typeVal)
				parseSnapshop(dec)
			}
			if typeVal == "l2update" {
				fmt.Println(typeVal)
				parseUpdate(dec)
			}

			fmt.Println(typeVal)
			fmt.Println("done")
		}

		closingBrace, err := dec.Token()
		if err != nil {
			log.Fatal(err)
		}
		if delim, ok := closingBrace.(json.Delim); !ok || delim != '}' {
			log.Println(closingBrace)
			log.Fatal("Expected object closing")
		}
	}

	// read closing bracket
	closingBracket, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	if delim, ok := closingBracket.(json.Delim); !ok || delim != ']' {
		log.Fatal("parse: Expected array closing")
	}
}

func parseUpdate(dec *json.Decoder) {
	for dec.More() {
		// parse changes
		k, err := dec.Token()
		if err != nil {
			log.Println(err)
		}
		// type key
		key, ok := k.(string)
		if !ok {
			continue
		}
		if key != "changes" {
			log.Println("expect changes")
			continue
		}
		fmt.Println(key)
		// changes
		parseChanges(dec)
	}
}

func parseChanges(dec *json.Decoder) {
	// array element
	openingBracket, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	if delim, ok := openingBracket.(json.Delim); !ok || delim != '[' {
		log.Fatal("Expected object")
	}

	for dec.More() {
		parseChange(dec)
	}

	// read closing bracket
	closingBracket, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	if delim, ok := closingBracket.(json.Delim); !ok || delim != ']' {
		log.Fatalf("parse array, expected array closing but got: %v", closingBracket)
	}
}

func parseChange(dec *json.Decoder) {
	openingBracket, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	if delim, ok := openingBracket.(json.Delim); !ok || delim != '[' {
		log.Fatal("Expected object")
	}

	side, err := dec.Token()
	if err != nil {
		log.Println(err)
	}
	amount, err := dec.Token()
	if err != nil {
		log.Println(err)
	}
	number, err := dec.Token()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(side, amount, number)

	closingBracket, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	if delim, ok := closingBracket.(json.Delim); !ok || delim != ']' {
		log.Fatal("Expected array closing")
	}
}

func parseSnapshop(dec *json.Decoder) {
	for dec.More() {
		// parse bids
		k, err := dec.Token()
		if err != nil {
			log.Println(err)
		}
		// type key
		key, ok := k.(string)
		if !ok {
			continue
		}
		if key != "bids" && key != "asks" {
			log.Println("expect bids or asks")
			continue
		}
		fmt.Println(key)
		// bids
		parseArray(dec)
	}
}

func parseArray(dec *json.Decoder) {
	// array element
	openingBracket, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	if delim, ok := openingBracket.(json.Delim); !ok || delim != '[' {
		log.Fatal("Expected object")
	}

	for dec.More() {
		parseValue(dec)
	}

	// read closing bracket
	closingBracket, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	if delim, ok := closingBracket.(json.Delim); !ok || delim != ']' {
		log.Fatalf("parse array, expected array closing but got: %v", closingBracket)
	}
}

func parseValue(dec *json.Decoder) {
	openingBracket, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	if delim, ok := openingBracket.(json.Delim); !ok || delim != '[' {
		log.Fatal("Expected object")
	}

	_, err = dec.Token()
	if err != nil {
		log.Println(err)
	}
	_, err = dec.Token()
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(k, v)

	closingBracket, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	if delim, ok := closingBracket.(json.Delim); !ok || delim != ']' {
		log.Fatal("Expected array closing")
	}
}
