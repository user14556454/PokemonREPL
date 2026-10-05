package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"project1/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
}

type config struct {
	commands    map[string]cliCommand
	nexturl     string
	previousurl string
	cached      *pokecache.Cache
	pokedex     map[string]pokemonCatch
}

func commandExit(c *config, p string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, p string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for cm := range c.commands {
		fmt.Printf("%v: %v\n", cm, c.commands[cm].description)
	}
	return nil
}

func commandMap(c *config, p string) error {
	var res *http.Response
	var err error
	var currenturl string
	if c.nexturl == "" {
		currenturl = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
		fmt.Println("map using api call")
		res, err = http.Get(currenturl)
		if err != nil {
			return err
		}
		defer res.Body.Close()
	} else {
		currenturl = c.nexturl
		res, err = http.Get(currenturl)
		if err != nil {
			return err
		}
		defer res.Body.Close()
	}
	var locname location
	jsondata, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	c.cached.Add(currenturl, jsondata)
	if err := json.Unmarshal(jsondata, &locname); err != nil {
		return err
	}
	for i := range locname.Results {
		fmt.Println(locname.Results[i].Name)
	}
	c.nexturl = locname.Next
	c.previousurl = locname.Previous
	return nil
}

func commandMapb(c *config, p string) error {
	if c.previousurl == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	var prevloc location
	cacheddata, rec := c.cached.Get(c.previousurl)
	if rec == true {
		fmt.Println("using cached data")
		if err := json.Unmarshal(cacheddata, &prevloc); err != nil {
			return err
		}
	} else {
		fmt.Println("api called")
		res, err := http.Get(c.previousurl)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		jsondata, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(jsondata, &prevloc); err != nil {
			return err
		}
	}
	for i := range prevloc.Results {
		fmt.Println(prevloc.Results[i].Name)
	}
	c.nexturl = prevloc.Next
	c.previousurl = prevloc.Previous
	return nil
}

func commandExplore(c *config, p string) error {
	cachedData, rec := c.cached.Get(p)
	var pokemonNames pokemonLocation
	if rec {
		fmt.Println("used cached data")
		if err := json.Unmarshal(cachedData, &pokemonNames); err != nil {
			return err
		}
	} else {
		fmt.Println("api call")
		reqURL := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", p)
		res, err := http.Get(reqURL)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		jsonData, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		c.cached.Add(p, jsonData)

		if err := json.Unmarshal(jsonData, &pokemonNames); err != nil {
			return err
		}
	}

	for _, pokem := range pokemonNames.PokemonEncounters {
		fmt.Println(pokem.Pokemon.PokeName)
	}
	return nil
}

func commandCatch(c *config, p string) error {
	fullURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", p)
	res, err := http.Get(fullURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var pokeCatch pokemonCatch
	if err := json.NewDecoder(res.Body).Decode(&pokeCatch); err != nil {
		return err
	}

	prob := (1.0 - (float32(pokeCatch.BaseExp) / 650.0)) * 100
	chance := rand.Intn(201)

	fmt.Printf("Throwing a Pokeball at %s...\n", p)
	if chance > int(prob) {
		fmt.Printf("%s escaped!\n", p)
		return nil
	}
	fmt.Printf("%s was caught!\n", p)
	fmt.Println("You may now inspect it with the inspect command.")
	c.pokedex[p] = pokeCatch
	return nil
}

func commandInspect(c *config, p string) error {
	for pokem := range c.pokedex {
		if pokem == p {
			fmt.Println("Name: ", p)
			fmt.Println("Height: ", c.pokedex[pokem].Height)
			fmt.Println("Weight: ", c.pokedex[pokem].Weight)
			fmt.Println("Stats:")
			for _, st := range c.pokedex[pokem].Stats {
				fmt.Printf("  -%s: %v\n", st.Stat.Name, st.BaseStat)
			}
			fmt.Println("Types:")
			for _, ty := range c.pokedex[pokem].Types {
				fmt.Println("  -", ty.Type.Name)
			}
			return nil
		}
	}
	fmt.Println("you have not caught that pokemon")
	return nil
}

func commandPokedex(c *config, p string) error {
	fmt.Println("Your Pokedex:")
	for pokem := range c.pokedex {
		fmt.Println("  -", pokem)
	}
	return nil
}

func runCommand(commandName string, commandParam string, c *config) {
	for cm := range c.commands {
		if cm == commandName {
			c.commands[commandName].callback(c, commandParam)
			return
		}
	}
	fmt.Println("Unknown Command")
}
