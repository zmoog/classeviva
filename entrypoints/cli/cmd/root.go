package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/adapters/feedback"
	"github.com/zmoog/classeviva/entrypoints/cli/cmd/agenda"
	"github.com/zmoog/classeviva/entrypoints/cli/cmd/grades"
)

var (
	debug    bool
	Feedback feedback.Feedback
)

func Execute() {
	var rootCmd = cobra.Command{
		PersistentPreRun: setupFeedback,
		Use:              "classeviva",
		Short:            "Classeviva is a CLI tool to access the popular school portal https://web.spaggiari.eu/",
	}

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Print debug information")

	rootCmd.AddCommand(agenda.NewCommand())
	rootCmd.AddCommand(grades.NewCommand())

	if err := rootCmd.Execute(); err != nil {
		Feedback.Error(err)
	}
}

func setupFeedback(cmd *cobra.Command, args []string) {
	Feedback = *feedback.Default()
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}
