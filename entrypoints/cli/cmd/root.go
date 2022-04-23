package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/entrypoints/cli/cmd/agenda"
	"github.com/zmoog/classeviva/entrypoints/cli/cmd/grades"
)

var (
	debug bool
)

func Execute() {
	var rootCmd = cobra.Command{
		PersistentPreRun: setupRootFlags,
		Use:              "classeviva",
		Short:            "Classeviva is a CLI tool to access the popular school portal https://web.spaggiari.eu/",
	}

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Print debug information")

	rootCmd.AddCommand(agenda.NewCommand())
	rootCmd.AddCommand(grades.NewCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func setupRootFlags(cmd *cobra.Command, args []string) {
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}
