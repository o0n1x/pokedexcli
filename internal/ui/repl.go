package ui

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	next     string
	previous string
}

var commands map[string]cliCommand
var defaultConfig config

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
	}
	defaultConfig = config{}
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

func ExecuteCommand(command string, cnfg *config) error {
	if cnfg == nil {
		cnfg = &defaultConfig
	}
	cmd, ok := commands[command]
	if ok {
		return cmd.callback(cnfg)
	} else {
		fmt.Println("Unknown command")
		return nil
	}

}

func commandExit(cnfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return errors.New("Failed to Exit")
}

func commandHelp(cnfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for cmd, disc := range commands {
		fmt.Printf("%s: %s\n", cmd, disc.description)
	}
	return nil
}

func commandMap(cnfg *config) error {
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

func commandMapb(cnfg *config) error {
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
