package main

import (
	"fmt"

	"github.com/twtiger/go-seccomp/parser"
)

func main() {
	result, err := parser.ParseFile("profiles/shared.seccomp")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Result: %#v\n", result)
	}
}
