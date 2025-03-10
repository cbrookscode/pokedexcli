package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name 		string
	description string
	callback 	func() error
}

func repl_loop() {
	scanner := bufio.NewScanner(os.Stdin)
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

		command, exists := getCommands()[commandname]
		if exists {
			err := command.callback()
			if err != nil {
				fmt.Errorf("Error executing callback for given command: %w", err)
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

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for key,value := range getCommands() {
		fmt.Printf("%s: %s\n", key, value.description)
	}
	fmt.Println()
	return nil
}