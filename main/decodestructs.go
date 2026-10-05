package main

type location struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []loactionArea `json:"results"`
}

type loactionArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type pokemonLocation struct {
	AreaName          string        `json:"name"`
	PokemonEncounters []pokemoninfo `json:"pokemon_encounters"`
}

type pokemoninfo struct {
	Pokemon poke `json:"pokemon"`
}

type poke struct {
	PokeName string `json:"name"`
	Pokeurl  string `json:"url"`
}

type pokemonCatch struct {
	Name    string  `json:"name"`
	Height  int     `json:"height"`
	Weight  int     `json:"weight"`
	BaseExp int     `json:"base_experience"`
	Stats   []stats `json:"stats"`
	Types   []types `json:"types"`
}

type stats struct {
	BaseStat int      `json:"base_stat"`
	Stat     statname `json:"stat"`
}

type statname struct {
	Name string `json:"name"`
}

type types struct {
	Type typename `json:"type"`
}

type typename struct {
	Name string `json:"name"`
}
