package main

import(
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	if text == "" {
		return []string{}
	}
	return strings.Fields(strings.ToLower(text))
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("Pokedex > ")
		
		// check if obtaining user input was successful. Clean input if so, else check for error. If user input wasn't nothing, print first word.
		var user_words []string
		if scanner.Scan() {
			user_words = cleanInput(scanner.Text())
		} else if err := scanner.Err(); err != nil {
			fmt.Errorf("There was an error: %w", err)
		}
		
		if len(user_words) != 0 {
			fmt.Printf("Your command was: %v\n", user_words[0])
		}
	}
}
