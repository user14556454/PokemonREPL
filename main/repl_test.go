package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "     How Are you    ",
			expected: []string{"how", "are", "you"},
		},
		{
			input:    "i am A Student",
			expected: []string{"i", "am", "a", "student"},
		},
	}
	for _, c := range cases {
		outputSlice := cleanInput(c.input)
		if len(outputSlice) != len(c.expected) {
			t.Errorf("Actual: %v, Expected: %v", len(outputSlice), len(c.expected))
		}
		for i := range outputSlice {
			actualword := outputSlice[i]
			expectedword := c.expected[i]
			if actualword != expectedword {
				t.Errorf("Actual: %v, Expected: %v", actualword, expectedword)
			}
		}
	}
}
