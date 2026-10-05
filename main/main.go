package main

import (
	"bufio"
	"fmt"
	"os"
	"project1/pokecache"
	"time"
)

func main() {
	sharedState := config{
		commands: map[string]cliCommand{
			"exit": {
				name:        "exit",
				description: "Exit the Pokedex",
				callback:    commandExit,
			},
			"help": {
				name:        "help",
				description: "Displays a help message",
				callback:    commandHelp,
			},
			"map": {
				name:        "map",
				description: "Gives location area",
				callback:    commandMap,
			},
			"mapb": {
				name:        "mapb",
				description: "Gives previous location area",
				callback:    commandMapb,
			},
			"explore": {
				name:        "explore",
				description: "Gives pokemons available in the location area",
				callback:    commandExplore,
			},
			"catch": {
				name:        "catch",
				description: "Catch a pokemon",
				callback:    commandCatch,
			},
			"inspect": {
				name:        "inspect",
				description: "Get details of a caught pokemon",
				callback:    commandInspect,
			},
			"pokedex": {
				name:        "pokedex",
				description: "Get names of all caught pokemons",
				callback:    commandPokedex,
			},
		},
		nexturl:     "",
		previousurl: "",
		cached:      pokecache.NewCache(5 * time.Second),
		pokedex:     map[string]pokemonCatch{},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		if scanner.Err() != nil {
			fmt.Println(scanner.Err())
			continue
		}
		inputSlice := cleanInput(scanner.Text())
		command := inputSlice[0]
		var param string
		if len(inputSlice) == 2 {
			param = inputSlice[1]
		}
		runCommand(command, param, &sharedState)
	}
}
