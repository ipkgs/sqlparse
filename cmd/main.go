package main

import (
	"fmt"
	"github.com/ipkgs/sqlparse"
	"os"
)

func main() {
	fmt.Println("Hello, World!")

	_, err := sqlparse.GetTokens(nil)

	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
