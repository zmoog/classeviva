package main

import (
	// Keep Europe/Rome calendar calculations available without system tzdata.
	_ "time/tzdata"

	"github.com/zmoog/classeviva/entrypoints/cli/cmd"
)

func main() {
	cmd.Execute()
}
