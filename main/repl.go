package main

import (
	"strings"
)

func cleanInput(text string) []string {
	words := []string{}
	text = strings.Trim(text, " ")
	count := 0
	word := ""
	for _, ch := range text {
		if ch == 32 {
			count++
			if count == 1 {
				words = append(words, word)
				word = ""
			}
			continue
		}
		count = 0
		word += string(ch)
	}
	if word != "" {
		words = append(words, word)
	}
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	return words
}
