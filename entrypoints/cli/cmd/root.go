package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/entrypoints/cli/cmd/grades"
)

func Execute() {
	var rootCmd = cobra.Command{
		Use:   "classeviva",
		Short: "Classeviva is a CLI tool to access the popular school portal https://web.spaggiari.eu/",
	}

	rootCmd.AddCommand(grades.NewCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
