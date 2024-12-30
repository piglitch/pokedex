package main

import "fmt"

func TestCleanInput(t *testing.T){
	cases := []struct {
		input string
		expect []string
	}{
		{
			input: "hello world",
			expect: []string{"hello", "world"},
		},
		{
			input: "I am playing",
			expect: []string{"I", "am", "playing"},
		},
		{
			input: "Boom! goes the dynamite",
			expect: []string{"Boom", "goes", "the", "dynamite"},
		},
	}
}