package ui

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/o0n1x/pokedexcli/internal/pokecache"
)

var cache pokecache.Cache

func init() {
	cache = *pokecache.NewCache(time.Second * 5)
}

type locationArea struct {
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []location `json:"results"`
}

type location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type locationinfo struct {
	PokemonEncounters []PokemonEncounters `json:"pokemon_encounters"`
}

type PokemonName struct {
	Name string `json:"name"`
}
type PokemonEncounters struct {
	Pokemon PokemonName `json:"pokemon"`
}

type Pokemon struct {
	BaseExperience int     `json:"base_experience"`
	Name           string  `json:"name"`
	Height         int     `json:"height"`
	Weight         int     `json:"weight"`
	Stats          []Stats `json:"stats"`
	Types          []Types `json:"types"`
}

type Stat struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Stats struct {
	BaseStat int  `json:"base_stat"`
	Effort   int  `json:"effort"`
	Stat     Stat `json:"stat"`
}
type Type struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Types struct {
	Slot int  `json:"slot"`
	Type Type `json:"type"`
}

func GetPokemon(pokemon string) (Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemon

	value, ok := cache.Get(url)
	if ok {
		poke := Pokemon{}
		err := json.Unmarshal(value, &poke)
		if err != nil {
			return Pokemon{}, err
		}
		return poke, nil
	}

	body, err := pokeclient(url)
	if err != nil {
		return Pokemon{}, err
	}
	cache.Add(url, body)
	poke := Pokemon{}
	err = json.Unmarshal(body, &poke)
	if err != nil {
		return Pokemon{}, err
	}
	return poke, nil

}

func GetLocationInfo(location string) (locationinfo, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + location

	value, ok := cache.Get(url)
	if ok {
		loc := locationinfo{}
		err := json.Unmarshal(value, &loc)
		if err != nil {
			return locationinfo{}, err
		}
		return loc, nil
	}

	body, err := pokeclient(url)
	if err != nil {
		return locationinfo{}, err
	}
	cache.Add(url, body)
	loc := locationinfo{}
	err = json.Unmarshal(body, &loc)
	if err != nil {
		return locationinfo{}, err
	}
	return loc, nil

}

func GetLocationArea(location_url string) (locationArea, error) {
	value, ok := cache.Get(location_url)
	if ok {
		loc := locationArea{}
		err := json.Unmarshal(value, &loc)
		if err != nil {
			return locationArea{}, err
		}
		return loc, nil
	}

	if location_url == "" {
		location_url = "https://pokeapi.co/api/v2/location-area"
	}
	body, err := pokeclient(location_url)
	if err != nil {
		return locationArea{}, err
	}
	cache.Add(location_url, body)
	loc := locationArea{}
	err = json.Unmarshal(body, &loc)
	if err != nil {
		return locationArea{}, err
	}
	return loc, nil
}

func pokeclient(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("ERROR: Response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	if err != nil {
		return nil, err
	}
	return body, nil
}
