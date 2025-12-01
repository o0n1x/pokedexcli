package ui

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type locationArea struct {
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []location `json:"results"`
}

type location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func GetLocationArea(location_url string) (locationArea, error) {
	if location_url == "" {
		location_url = "https://pokeapi.co/api/v2/location-area"
	}
	body, err := pokeclient(location_url)
	if err != nil {
		return locationArea{}, err
	}
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
