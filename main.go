package main

import (
	"fmt"

	seccomp "github.com/subgraph/go-seccomp"
)

func main() {
	result, _ := seccomp.Compile("bla.seccomp", true)
	for _, filter := range result {
		fmt.Println(filter)
	}

	// result, err := parser.ParseFile("profiles/shared.seccomp")
	// if err != nil {
	// 	fmt.Printf("Error: %s\n", err)
	// } else {
	// 	fmt.Printf("Result: %#v\n", result)
	// }
}
