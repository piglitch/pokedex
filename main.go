package main

import (
	"bufio"
	"math/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/piglitch/pokedexcli/pokecache"
)

type Pokemon struct {
	name string
	height int
	weight int
	hp int
	attack int
	defense int
	specialAttack int
	specialDefense int
	speed int
	types []string
}

type PokemonResponseByName struct {
	Name string `json:"name"`
	BaseExp int `json:"base_experience"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct{
		PokeType struct{
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
	Height int `json:"height"`
}

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
	pokemons := make(map[string]Pokemon)
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
				return commandMap(conf, P, userInputSlice, pokemons)
			},
		},
		"mapb": {
			name: "mapb",
			description: "Shows previous 20 location areas in the pokemon world.",
			callback: func(conf *Config, userInputSlice []string) error{
				return commandMapb(conf, P, userInputSlice, pokemons)
			},
		},
		"explore": {
			name: "explore",
			description: "Shows pokemons of a certain location area.",
			callback: func(conf *Config, userInputSlice []string) error{
				return commandExplore(conf, P, userInputSlice, pokemons)
			},
		},
		"catch": {
			name: "catch",
			description: "Attempts to catch a pokemon",
			callback: func (conf *Config, userInputSlice []string) error {
				return commandCatch(conf, P, userInputSlice, pokemons)
			},
		},
		"inspect": {
			name: "inspect",
			description: "Inspects a pokemon once caught",
			callback: func(conf *Config, userInputSlice []string) error {
				return commandInspect(conf, P, userInputSlice, pokemons)
			},
		},
		"pokedex": {
			name: "pokedex",
			description: "Shows the pokemons you have caught",
			callback: func(conf *Config, userInputSlice []string) error {
				return commandPokedex(conf, P, userInputSlice, pokemons)
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

func commandHelp(conf *Config, userInputSlice []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Available commands:")

	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("map: Displays the names of 20 location areas in the Pokemon world.")
	fmt.Println("mapb: Shows previous 20 location areas in the pokemon world.")
	fmt.Println("explore <location>: Shows pokemons of a certain location area.")
	fmt.Println("catch <pokemon>: Attempts to catch a pokemon")
	fmt.Println("inspect <pokemon>: Inspects a pokemon once caught")
	fmt.Println("pokedex: Shows the pokemons you have caught")

	return nil
}


func commandMap(conf *Config, P *pokecache.Cache, _ []string, _ map[string]Pokemon) error {
	
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

func commandMapb(conf *Config, P *pokecache.Cache, _ []string, _ map[string]Pokemon) error {

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

func commandExplore(_ *Config, P *pokecache.Cache, userInputSlice []string, _ map[string]Pokemon) error {

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

func commandCatch(_ *Config, _ *pokecache.Cache, userInputSlice []string, pokemons map[string]Pokemon) error {
	_, exists := pokemons[userInputSlice[1]]
	if exists {
		fmt.Printf("You have already caught %s \n", userInputSlice[1])
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %s... \n", userInputSlice[1])
	pokemonUrl := "https://pokeapi.co/api/v2/pokemon/"
	fulUrl := pokemonUrl + userInputSlice[1]

	res, err := http.Get(fulUrl)
	if err != nil {
		return fmt.Errorf("failed fetching pokemon: %s", err)
	}
	defer res.Body.Close()
	var data PokemonResponseByName
	body, err := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		fmt.Printf("%s is not a pokemon: %d\n", userInputSlice[1], res.StatusCode)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading response body: %s", err)
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body: %s", err)
	}
	pokeExp := data.BaseExp
	chances := rand.Float64() * 2

	var hp int
	var attack int
	var defense int
	var specialAttack int
	var specialDefense int
	var speed int
	pokemonTypes := []string{}
	for _, stats := range data.Stats {
		if stats.Stat.Name == "hp" {
			hp = stats.BaseStat
		}
		if stats.Stat.Name == "attack" {
			attack = stats.BaseStat
		}
		if stats.Stat.Name == "defense" {
			defense = stats.BaseStat
		}
		if stats.Stat.Name == "special-attack" {
			specialAttack = stats.BaseStat
		}
		if stats.Stat.Name == "special-defense" {
			specialDefense = stats.BaseStat
		}
		if stats.Stat.Name == "speed" {
			speed = stats.BaseStat
		}
	}

	for _, types := range data.Types {
		pokemonTypes = append(pokemonTypes, types.PokeType.Name)
	}

	if chances * 100 >= float64(pokeExp) {
		fmt.Printf("%s was caught \n", userInputSlice[1])
		pokemons[userInputSlice[1]] = Pokemon{
			name: data.Name,
			height: data.Height,
			weight: data.Weight,
			hp: hp,
			attack: attack,
			defense: defense,
			specialAttack: specialAttack,
			specialDefense: specialDefense,
			speed: speed,
			types: pokemonTypes,
		}
	} else {
		fmt.Printf("%s escaped \n", userInputSlice[1])
	}
	fmt.Printf("caught pokemons: %v \n", pokemons)
	return nil
}

func commandInspect (_ *Config, _ *pokecache.Cache, userInputSlice []string, pokemons map[string]Pokemon) error {
	pokemon, exists := pokemons[userInputSlice[1]]
	if !exists {
		fmt.Printf("You haven't caught %s. You have to catch a pokemon to inspect it. \n", userInputSlice[1])
		return nil
	}
	fmt.Printf("Name: %s\n", pokemon.name)
	fmt.Printf("Height: %d\n", pokemon.height)
	fmt.Printf("Weight: %d\n", pokemon.weight)
	fmt.Println("Stats: ")
	fmt.Printf("	-hp: %d\n", pokemon.hp)
	fmt.Printf("	-attack: %d\n", pokemon.attack)
	fmt.Printf("	-defense: %d\n", pokemon.defense)
	fmt.Printf("	-special-attack: %d\n", pokemon.specialAttack)
	fmt.Printf("	-special-defense: %d\n", pokemon.specialDefense)
	fmt.Printf("	-speed: %d\n", pokemon.speed)
	fmt.Println("Types: ")
	
	for _, pokeType := range pokemon.types {
		fmt.Printf("	- %s \n", pokeType)	
	}

	return nil
}

func commandPokedex (_ *Config, _ *pokecache.Cache, _ []string, pokemons map[string]Pokemon) error {
	if len(pokemons) < 1 {
		fmt.Println("You have not caught any pokemons yet")
	}
	for _, pokemon := range pokemons {
		fmt.Printf("	- %s \n", pokemon.name)
	}
	return nil
}