package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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
	Callback    func(*Config, ...string) error
}

type LocationData struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int           `json:"chance"`
				ConditionValues []interface{} `json:"condition_values"`
				MaxLevel        int           `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
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
		"explore": {
			Name:        "explore",
			Description: "Displays pokemon names for provided location",
			Callback:    CommandExplore,
		},
	}
}

func CommandExit(c *Config, params ...string) error {
	_ = c
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp(c *Config, params ...string) error {
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
func getBytesFromHttpGet(endpoint_url string) ([]byte, error) {
	res, err := http.Get(endpoint_url)
	if err != nil {
		return []byte{}, fmt.Errorf("error with get request: error = %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return []byte{}, fmt.Errorf("error with response from get request: status = %d", res.StatusCode)
	}

	raw_bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("error reading response body using io.ReadAll: error = %w", err)
	}
	return raw_bytes, nil
}

// Get raw []byte data and store in cache
func getPokedata(endpoint_url string) error {
	raw_bytes, err := getBytesFromHttpGet(endpoint_url)
	if err != nil {
		return fmt.Errorf("error getting bytes from http")
	}
	cache.Add(endpoint_url, raw_bytes)

	return nil
}

func structureDataIntoPointerStruct[T any](endpoint_url string, c *T) error {
	_, exist := cache.Get(endpoint_url)
	if !exist {
		fmt.Println("making get request to retrieve data")
		err := getPokedata(endpoint_url)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("data exists in cache")
	}

	data, _ := cache.Get(endpoint_url)
	err := json.Unmarshal(data, c)
	if err != nil {
		return fmt.Errorf("error unmarshalling json")
	}
	return nil
}

func CommandMap(c *Config, params ...string) error {
	endpoint_url := "https://pokeapi.co/api/v2/location-area/"

	if c.Next != nil {
		endpoint_url = *c.Next
	}

	err := structureDataIntoPointerStruct(endpoint_url, c)
	if err != nil {
		return err
	}

	for _, location := range c.Results {
		fmt.Printf("location name: %v\n", location.Name)
	}

	return nil
}

func CommandMapb(c *Config, params ...string) error {
	if c.Previous == nil {
		return errors.New("previous page doesn't exist")
	}

	endpoint_url := *c.Previous

	err := structureDataIntoPointerStruct(endpoint_url, c)
	if err != nil {
		return err
	}

	for _, location := range c.Results {
		fmt.Printf("location name: %v\n", location.Name)
	}

	return nil
}

// seperate func to retrieve data from cache
func CommandExplore(c *Config, params ...string) error {
	_ = c

	loc_data := &LocationData{}

	if len(params) == 0 {
		return fmt.Errorf("need to provide location name")
	}
	loc_name := strings.Join(params, "-")
	key := "https://pokeapi.co/api/v2/location-area/" + loc_name

	err := structureDataIntoPointerStruct(key, loc_data)
	if err != nil {
		return fmt.Errorf("invalid location name")
	}

	fmt.Printf("Exploring %s\n", loc_name)
	fmt.Println("Found Pokemon:")
	for _, pokemon_struct := range loc_data.PokemonEncounters {
		fmt.Printf("- %s\n", pokemon_struct.Pokemon.Name)
	}

	return nil
}
