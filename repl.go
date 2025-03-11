package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	commands "github.com/cbrookscode/pokedexcli/internal"
)

func repl_loop() {
	scanner := bufio.NewScanner(os.Stdin)
	Config_pointer := &commands.Config{}

	for {
		fmt.Printf("Pokedex > ")

		// check if obtaining user input was successful. Clean input if so, else check for error.
		var user_words []string
		if scanner.Scan() {
			user_words = cleanInput(scanner.Text())
		} else if err := scanner.Err(); err != nil {
			fmt.Errorf("There was an error: %w", err)
		}

		// if no input skip to next iteration of loop
		if len(user_words) == 0 {
			continue
		}

		// if command exists in supported commands, call its callback function. if there is an error, print it. if command doesn't exist let user know.
		commandname := user_words[0]

		command, exists := commands.GetCommands()[commandname]
		if exists {
			err := command.Callback(Config_pointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func cleanInput(text string) []string {
	if text == "" {
		return []string{}
	}
	return strings.Fields(strings.ToLower(text))
}
