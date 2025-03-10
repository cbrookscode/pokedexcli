package main

import(
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	if text == "" {
		return []string{}
	}
	return strings.Fields(strings.ToLower(text))
}

func main() {
	fmt.Println("Hello, World!")
}
