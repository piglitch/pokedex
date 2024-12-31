package main

import "testing"

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
			expect: []string{"Boom!", "goes", "the", "dynamite"},
		},
	}

	for _, c := range cases{
		actual := cleanInput(c.input)

		for i := range actual{
			word := actual[i]
			expectedWord := c.expect[i]

			if word != expectedWord {
				t.Errorf("Test failed, words don't match; Actual: %s. Wanted: %s.", word, expectedWord)
				return
			}
		}
	}
}