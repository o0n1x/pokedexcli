package ui

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type config struct {
	next     string
	previous string
}

var commands map[string]cliCommand
var defaultConfig config
var pokedex map[string]Pokemon

func init() {
	commands = map[string]cliCommand{
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
			description: "Displays the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations if any",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "explore a given location and returns all pokemons in that location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "attempt to catch the pokemon given",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "inspect a caught pokemon",
			callback:    commandInspect,
		},
	}
	defaultConfig = config{}
	pokedex = make(map[string]Pokemon)
}

func CleanInput(text string) []string {
	rslt := strings.Split(strings.ToLower(text), " ")

	var output []string

	for _, str := range rslt {
		if str != "" && str != " " {
			output = append(output, str)
		}
	}
	return output
}

func ExecuteCommand(command []string, cnfg *config) error {
	if cnfg == nil {
		cnfg = &defaultConfig
	}
	cmd, ok := commands[command[0]]
	var args = []string{""}
	if len(command) > 1 {
		args = command[1:]
	}
	if ok {
		return cmd.callback(cnfg, args)
	} else {
		fmt.Println("Unknown command")
		return nil
	}

}

func commandExit(_ *config, _ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return errors.New("Failed to Exit")
}

func commandHelp(_ *config, _ []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for cmd, disc := range commands {
		fmt.Printf("%s: %s\n", cmd, disc.description)
	}
	return nil
}

func commandMap(cnfg *config, _ []string) error {
	rslt, err := GetLocationArea(cnfg.next)
	if err != nil {
		return err
	}
	cnfg.next = rslt.Next
	cnfg.previous = rslt.Previous
	locations := rslt.Results

	for _, loc := range locations {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(cnfg *config, _ []string) error {
	rslt, err := GetLocationArea(cnfg.previous)
	if err != nil {
		return err
	}
	cnfg.next = rslt.Next
	cnfg.previous = rslt.Previous
	locations := rslt.Results

	for _, loc := range locations {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandExplore(_ *config, args []string) error {

	location := args[0]
	fmt.Printf("Exploring %s...\n", location)
	rslt, err := GetLocationInfo(location)
	if err != nil {
		return err
	}
	for _, pokemon := range rslt.PokemonEncounters {
		fmt.Println("  -", pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(_ *config, args []string) error {
	pokemon := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)
	rslt, err := GetPokemon(pokemon)
	if err != nil {
		return err
	}

	chance := 50.0 / float64(rslt.BaseExperience)
	if chance > rand.Float64() {
		fmt.Printf("%s was caught!\n", pokemon)
		pokedex[pokemon] = rslt
	} else {
		fmt.Printf("%s escaped!\n", pokemon)
	}
	return nil
}

func commandInspect(_ *config, args []string) error {

	pokemon, ok := pokedex[args[0]]
	if !ok {
		fmt.Println("Pokemon not caught!")
		return nil
	}
	fmt.Println("Name:", pokemon.Name)
	fmt.Println("Height:", pokemon.Height)
	fmt.Println("Weight:", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("	-%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, types := range pokemon.Types {
		fmt.Printf("	- %s\n", types.Type.Name)
	}
	return nil
}
