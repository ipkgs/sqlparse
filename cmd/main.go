package main

import (
	"fmt"
	"github.com/ipkgs/sqlparse"
	"os"
	"strings"
)

func main() {
	fmt.Println("sqlparse")

	tokens, err := sqlparse.GetTokens(strings.Join(os.Args[1:], " "))

	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	for _, token := range tokens {
		fmt.Printf("%s: %s\n", token.Type, token.Value)
	}
}
