package main

import (
	"fmt"
	"strings"
)

func main(){
	// fmt.Println("Hello, World!")	
	fmt.Println(cleanInput("messi is not the goat"))
}

func cleanInput(text string) []string{
	return strings.Split(text, " ")
}