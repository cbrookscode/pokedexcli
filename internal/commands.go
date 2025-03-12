package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var cache *Cache

func init() {
	cache = NewCache(1200 * time.Second)
}

type Config struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*Config) error
}

func GetCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    CommandExit,
		},
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    CommandHelp,
		},
		"map": {
			Name:        "map",
			Description: "Displays 20 locations",
			Callback:    CommandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Displays previous 20 locations",
			Callback:    CommandMapb,
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
	for key, value := range GetCommands() {
		fmt.Printf("%s: %s\n", key, value.Description)
	}
	fmt.Println()
	return nil
}

func getLocations(c *Config, thecache *Cache, endpoint_url string) error {

	value, exists := thecache.Get(endpoint_url)

	if exists {
		// using cache
		err := json.Unmarshal(value, c)
		if err != nil {
			return fmt.Errorf("issue with unmarshaling json from cache: %w", err)
		}
		for _, location := range c.Results {
			fmt.Printf("location name: %v\n", location.Name)
		}
	} else {
		// making get request
		res, err := http.Get(endpoint_url)
		if err != nil {
			return fmt.Errorf("error with get request: error = %w", err)
		}
		defer res.Body.Close()
		if res.StatusCode < 200 && res.StatusCode > 299 {
			return fmt.Errorf("error with response from get request: status = %d", res.StatusCode)
		}

		raw_bytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading response body using io.ReadAll: error = %w", err)
		}

		err = json.Unmarshal(raw_bytes, c)
		if err != nil {
			return fmt.Errorf("error unmarshalling from new get request: error = %w", err)
		}

		for _, location := range c.Results {
			fmt.Printf("location name: %v\n", location.Name)
		}

		thecache.Add(endpoint_url, raw_bytes)
	}
	return nil
}

func CommandMap(c *Config) error {
	endpoint_url := "https://pokeapi.co/api/v2/location-area/"

	if c.Next != nil {
		endpoint_url = *c.Next
	}
	err := getLocations(c, cache, endpoint_url)

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

	err := getLocations(c, cache, endpoint_url)

	if err != nil {
		return err
	}

	return nil
}
