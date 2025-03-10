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

var supported_cmds map[string]cliCommand

func init() {
	supported_cmds = map[string]cliCommand{
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
		command := user_words[0]
		if _,ok := supported_cmds[command]; ok {
			err := supported_cmds[command].callback()
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

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for key,value := range supported_cmds {
		fmt.Printf("%s: %s\n", key, value.description)
	}
	return nil
}