package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type Config struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type CliCommand struct {
	Name 		string
	Description string
	Callback 	func(*Config) error
}

func GetCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"exit": {
			Name: "exit",
			Description: "Exit the Pokedex",
			Callback: CommandExit,
		},
		"help": {
			Name: "help",
			Description: "Displays a help message",
			Callback: CommandHelp,
		},
		"map": {
			Name: "map",
			Description: "Displays 20 locations",
			Callback: CommandMap,
		},
		"mapb": {
			Name: "mapb",
			Description: "Displays previous 20 locations",
			Callback: CommandMapb,
		},
	}
}

func CommandExit(c *Config) error {
	_ = c
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp(c *Config) error {
	_ = c
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for key,value := range GetCommands() {
		fmt.Printf("%s: %s\n", key, value.Description)
	}
	fmt.Println()
	return nil
}

func getLocations(c *Config, endpoint_url string) error {
	res, err := http.Get(endpoint_url)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 && res.StatusCode > 299 {
		return fmt.Errorf("Error with response from Get request: Status = %v", res.StatusCode)
	}
	defer res.Body.Close()


	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(c)
	if err != nil {
		return err
	}

	for _,location := range c.Results {
		fmt.Printf("location name: %v\n", location.Name)
	}
	return nil
}

func CommandMap(c *Config) error {
	endpoint_url := "https://pokeapi.co/api/v2/location-area/"

	if c.Next != nil {
		endpoint_url = *c.Next
	}
	err := getLocations(c, endpoint_url)

	if err != nil {
		return err
	}
	
	return nil
}

func CommandMapb(c *Config) error {

	if c.Previous == nil {
		return errors.New("Previous page doesn't exist")
	}

	endpoint_url := *c.Previous

	err := getLocations(c, endpoint_url)

	if err != nil {
		return err
	}

	return nil
}