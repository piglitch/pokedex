package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name string
	description string
	callback func() error
}

type config struct {
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
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Pokedex > ")
	commands := map[string]cliCommand{
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
			callback: commandMap,
		},
	}
	for scanner.Scan() {
		userInput := scanner.Text()
		userInput = strings.ToLower(userInput)
		found := false
		mapCalled := false
		for _, cmd := range commands {
			if cmd.name == userInput {
				found = true
				cmd.callback()
			}
			if userInput == "map" {
				mapCalled = true
			}
		}
		if !found {
			fmt.Println("Unknown command")
		}
	}
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye! \n")
	os.Exit(0)
	return errors.New("whatever")
}

func commandHelp() error{
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	return errors.New("whatever")
}

func commandMap() error{
	pokeUrl := "https://pokeapi.co/api/v2/location-area/?offset=20"
	res, err := http.Get(pokeUrl)
	if err != nil {
		return fmt.Errorf("unable to fetch!: %s", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != err {
		return fmt.Errorf("could not read response body: %s", err)
	}
	var data LocationAreaResponse
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %s", err)
	}
	locations := data.Results
	for _, loc := range locations{
		fmt.Println(loc.Name)
	}
	return fmt.Errorf("whatever")
}

func cleanInput(text string) []string{
	return strings.Split(text, " ")
}