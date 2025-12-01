package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/o0n1x/pokedexcli/internal/ui"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		if scanner.Err() != nil {
			fmt.Printf("Error: %v \n", scanner.Err())
		}

		txt := scanner.Text()

		cleantxt := ui.CleanInput(txt)[0]

		ui.ExecuteCommand(cleantxt, nil)
	}
}
