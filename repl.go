package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	commands "github.com/cbrookscode/pokedexcli/internal"
)

func repl_loop() {
	Config_pointer := &commands.Config{}

	for {
		fmt.Printf("Pokedex > ")

		// check if obtaining user input was successful. Clean input if so, else check for error.
		user_words, err := getUserInput()
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		// if no input skip to next iteration of loop
		if len(user_words) == 0 {
			continue
		}

		// if command exists in supported commands, call its callback function. if there is an error, print it. if command doesn't exist let user know.
		commandname := user_words[0]
		params := user_words[1:]

		command, exists := commands.GetCommands()[commandname]
		if exists {
			err = command.Callback(Config_pointer, params...)
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

func getUserInput() ([]string, error) {
	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() {
		user_words := cleanInput(scanner.Text())
		return user_words, nil
	} else if err := scanner.Err(); err != nil {
		return []string{}, fmt.Errorf("error getting user input: %w", err)
	} else {
		return []string{}, io.EOF
	}
}
