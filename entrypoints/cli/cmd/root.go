package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/adapters/feedback"
	"github.com/zmoog/classeviva/entrypoints/cli/cmd/agenda"
	"github.com/zmoog/classeviva/entrypoints/cli/cmd/grades"
)

var (
	debug  bool
	format string
)

func Execute() {
	var rootCmd = cobra.Command{
		PersistentPreRun: configureFeedback,
		Use:              "classeviva",
		Short:            "Classeviva is a CLI tool to access the popular school portal https://web.spaggiari.eu/",
	}

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Print debug information")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "text", "Output format")

	rootCmd.AddCommand(agenda.NewCommand())
	rootCmd.AddCommand(grades.NewCommand())

	if err := rootCmd.Execute(); err != nil {
		feedback.Error(err)
	}
}

func configureFeedback(cmd *cobra.Command, args []string) {
	switch format {
	case "json":
		feedback.SetFormat(feedback.JSON)
	default:
		feedback.SetFormat(feedback.Text)
	}

	if debug {
		log.SetLevel(log.TraceLevel)
	}
}
