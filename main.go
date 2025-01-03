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
	callback func(*Config) error
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
			callback:	func(conf *Config) error{
				return commandMap(conf, P)
			},
		},
		"mapb": {
			name: "mapb",
			description: "Shows previous 20 location areas in the pokemon world.",
			callback: func(conf *Config) error{
				return commandMapb(conf, P)
			},
		},
		"explore": {
			name: "explore",
			description: "Shows pokemons of a certain location area.",
			callback: func(conf *Config) error{
				return commandMapb(conf, P)
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
		found := false

		for _, cmd := range commands {
			if cmd.name == userInput {
				found = true
				cmd.callback(&preConf)
			}
		}
		if !found {
			fmt.Println("Unknown command")
		}
	}
}

func commandExit(conf *Config) error {
	fmt.Print("Closing the Pokedex... Goodbye! \n")
	os.Exit(0)
	return nil
}

func commandHelp(conf *Config) error{
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}

func commandMap(conf *Config, P *pokecache.Cache) error{
	
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
	
	for _, loc := range locations{
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(conf *Config, P *pokecache.Cache) error {

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

func commandExplore(C *Config, P *pokecache.Cache) error {
	exploreUrl := "https://pokeapi.co/api/v2/location/" + 
	res, err := http.Get() 	
}
