package main

import (
	"fmt"
	"log"

	"github.com/masp/hoser/lexer"
)

func main() {
	tokens, err := lexer.ScanAll("I have a new toy\n	it is called a lexer and it is great")
	fmt.Printf("%v\n", tokens)
	if err != nil {
		log.Fatalf("failed to scan: %v", err)
	}
}
