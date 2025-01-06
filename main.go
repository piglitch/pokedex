package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/piglitch/pokedexcli/pokecache"
)



type CliCommand struct {
	name string
	description string
	callback func(*Config, []string) error
}

type Config struct {
	Next string
	Previous string
}

type LocationAreaResponse struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type PokemonResponse struct {
	PokeEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func main(){
	P := pokecache.NewCache(120 * time.Second) 
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Pokedex > ")
	commands := map[string]CliCommand{
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
		"map": {
			name: "map",
			description: "Displays the names of 20 location areas in the Pokemon world.",
			callback:	func(conf *Config, userInputSlice []string) error{
				return commandMap(conf, P, userInputSlice)
			},
		},
		"mapb": {
			name: "mapb",
			description: "Shows previous 20 location areas in the pokemon world.",
			callback: func(conf *Config, userInputSlice []string) error{
				return commandMapb(conf, P, userInputSlice)
			},
		},
		"explore": {
			name: "explore",
			description: "Shows pokemons of a certain location area.",
			callback: func(conf *Config, userInputSlice []string) error{
				return commandExplore(conf, P, userInputSlice)
			},
		},
		"catch": {
			name: "catch",
			description: "Attempts to catch a pokemon",
			callback: func (conf *Config, userInputSlice []string) error {
				return commandCatch(conf, P, userInputSlice)
			},
		},
	}
	preConf := Config{
		Next: "",
		Previous: "",
	}
	for scanner.Scan() {
		userInput := scanner.Text()
		userInput = strings.ToLower(userInput)
		userInputSlice := strings.SplitN(userInput, " ", 2)
		found := false

		for _, cmd := range commands {
			if cmd.name == userInputSlice[0] {
				found = true
				cmd.callback(&preConf, userInputSlice)
			}
		}
		if !found {
			fmt.Println("Unknown command")
		}
	}
}

func commandExit(conf *Config, userInputSlice []string) error {
	fmt.Print("Closing the Pokedex... Goodbye! \n")
	os.Exit(0)
	return nil
}

func commandHelp(conf *Config, userInputSlice []string) error{
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}

func commandMap(conf *Config, P *pokecache.Cache, userInputSlice []string) error {
	
	pokeUrl := conf.Next

	if pokeUrl == "" {
		pokeUrl = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"	
	}
	var data LocationAreaResponse
	entry, exists := P.Get(pokeUrl)
	if exists {
		err := json.Unmarshal([]byte(entry), &data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal: %s", err)
		}
	} else {
		res, err := http.Get(pokeUrl)
		if err != nil {
			return fmt.Errorf("unable to fetch!: %s", err)
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)

		if err != nil {
			return fmt.Errorf("could not read response body: %s", err)
		}

		err = json.Unmarshal([]byte(body), &data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal: %s", err)
		}
		P.Add(pokeUrl, body)	
	}

	conf.Next = data.Next
	conf.Previous = data.Previous
	
	locations := data.Results
	
	for _, loc := range locations {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(conf *Config, P *pokecache.Cache, userInputSlice []string) error {

	var data LocationAreaResponse
	pokeUrl := conf.Previous
	if pokeUrl == "" {
		fmt.Println("you're on the first page")	
	}
	entry, exists := P.Get(pokeUrl)
	if exists {
		err := json.Unmarshal([]byte(entry), &data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal: %s", err)
		}
	 } 
	conf.Next = data.Next
	conf.Previous = data.Previous

	locations := data.Results
	
	for _, loc := range locations{
		fmt.Println(loc.Name)
	}
	return nil
}

func commandExplore(C *Config, P *pokecache.Cache, userInputSlice []string) error {

	var data PokemonResponse
	locUrl := "https://pokeapi.co/api/v2/location-area/" 
	fullUrl := locUrl + userInputSlice[1]

	entry, exists := P.Get(fullUrl)
	if exists {
		err := json.Unmarshal(entry, &data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal: %s", err)
		}

	} else {
		res, err := http.Get(fullUrl) 	
		if err != nil {
			return fmt.Errorf("failed to fetch %s", err)
		}
		defer res.Body.Close()
	
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed reading response body: %s", err)
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal response body: %s", err)
		}
		P.Add(fullUrl, body)
	}
	
	for _, encounter := range data.PokeEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(C *Config, P *pokecache.Cache, userInputSlice []string) error {
	
	return nil
}