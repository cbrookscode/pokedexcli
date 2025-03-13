package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var cache *Cache
var Pokedex map[string]Pokemon

func init() {
	cache = NewCache(1200 * time.Second)
	Pokedex = make(map[string]Pokemon)
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
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
	BaseExperience         int    `json:"base_experience"`
	Height                 int    `json:"height"`
	ID                     int    `json:"id"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name    string `json:"name"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
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
		"catch": {
			Name:        "catch",
			Description: "Try to catch a pokemon and add it to your pokedex",
			Callback:    CommandCatch,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Get the details of a Pokemon you have caught",
			Callback:    CommandInspect,
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

func calcChanceToCatch(basexp int) int {
	chance_to_catch := 100 - (5 * (basexp / 30))
	if chance_to_catch > 95 {
		chance_to_catch = 95
	} else if chance_to_catch < 5 {
		chance_to_catch = 5
	}
	return chance_to_catch
}

func CommandCatch(c *Config, params ...string) error {
	poke_struct := &Pokemon{}

	if len(params) == 0 {
		return fmt.Errorf("need to provide pokemon name")
	}
	pokemon_name := strings.Join(params, "-")
	endpoint_url := "https://pokeapi.co/api/v2/pokemon/" + pokemon_name

	data, err := getBytesFromHttpGet(endpoint_url)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, poke_struct)
	if err != nil {
		return err
	}
	chance := calcChanceToCatch(poke_struct.BaseExperience)
	rnum := rand.Intn(100)
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon_name)
	if rnum > chance {
		fmt.Printf("%s escaped!\n", pokemon_name)
		return nil
	} else {
		fmt.Printf("%s was caught!\n", pokemon_name)
		Pokedex[pokemon_name] = *poke_struct
		return nil
	}
}

func CommandInspect(c *Config, params ...string) error {
	if len(params) == 0 {
		return fmt.Errorf("need to provide pokemon name")
	}
	pokemon_name := strings.Join(params, "-")

	poke_struct, exists := Pokedex[pokemon_name]
	if exists {
		fmt.Printf("Name: %s\n", poke_struct.Name)
		fmt.Printf("Height: %d\n", poke_struct.Height)
		fmt.Printf("Weight: %d\n", poke_struct.Weight)
		fmt.Println("Stats:")
		for _, stat := range poke_struct.Stats {
			fmt.Printf("   -%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, poketype := range poke_struct.Types {
			fmt.Printf("   - %s\n", poketype.Type.Name)
		}
		return nil
	}
	fmt.Println("you have not caught that pokemon")
	return nil
}
