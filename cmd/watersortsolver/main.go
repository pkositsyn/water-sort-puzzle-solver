package main

import (
	"fmt"

	watersortpuzzle "github.com/pkositsyn/water-sort-puzzle-solver"
)

func main() {
	fmt.Println("Input initial puzzle state")

	var initialStateStr string
	n, err := fmt.Scanln(&initialStateStr)
	if err != nil {
		fmt.Printf("Error getting input: %s\n", err.Error())
		return
	}
	if n != 1 {
		fmt.Printf("Scanned %d values but needed one position\n", n)
		return
	}

	solver := watersortpuzzle.NewAStarSolver()

	var initialState watersortpuzzle.State
	if err := initialState.FromString(initialStateStr); err != nil {
		fmt.Printf("Invalid puzzle state provided: %s\n", err.Error())
		return
	}

	steps, err := solver.Solve(initialState)
	if err != nil {
		fmt.Printf("Cannot solve puzzle: %s\n", err.Error())
		return
	}

	fmt.Printf("Puzzle solved in %d steps!\n", len(steps))
	for _, step := range steps {
		fmt.Println(step.From+1, step.To+1)
	}
}
