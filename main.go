package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main(){
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Pokedex >")
	for scanner.Scan() {
		userInput := scanner.Text()
		userInput = strings.ToLower(userInput)
		words := strings.Fields(userInput)
		fmt.Printf("Your command was: %s \n", words[0])
		fmt.Print("Pokedex >")
	}
}

func cleanInput(text string) []string{
	return strings.Split(text, " ")
}