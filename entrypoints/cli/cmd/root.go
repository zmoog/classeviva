package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/adapters/feedback"
	"github.com/zmoog/classeviva/commands"
	"github.com/zmoog/classeviva/entrypoints/cli/cmd/agenda"
	"github.com/zmoog/classeviva/entrypoints/cli/cmd/grades"
	"github.com/zmoog/classeviva/entrypoints/cli/cmd/noticeboards"
	"github.com/zmoog/classeviva/entrypoints/cli/cmd/profile"
	"github.com/zmoog/classeviva/entrypoints/cli/cmd/version"
)

var (
	debug       bool
	format      string
	profileFlag string
	username    string
	password    string
)

// GetRunnerOpts is a function variable that subcommands can call to get runner options
// It's set after flag parsing to avoid import cycles
var GetRunnerOpts func() commands.RunnerOptions

func Execute() {
	// Set up the function to get runner options (avoids import cycle)
	GetRunnerOpts = func() commands.RunnerOptions {
		return commands.RunnerOptions{
			Username: username,
			Password: password,
			Profile:  profileFlag,
		}
	}

	var rootCmd = cobra.Command{
		PersistentPreRun: configureFeedback,
		Use:              "classeviva",
		Short:            "Classeviva is a CLI tool to access the popular school portal https://web.spaggiari.eu/",
	}

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Print debug information")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "text", "Output format")
	rootCmd.PersistentFlags().StringVarP(&profileFlag, "profile", "p", "", "Profile name to use")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "Username (overrides profile)")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Password (overrides profile)")

	rootCmd.AddCommand(agenda.NewCommand())
	rootCmd.AddCommand(grades.NewCommand())
	rootCmd.AddCommand(noticeboards.NewCommand())
	rootCmd.AddCommand(profile.NewCommand())
	rootCmd.AddCommand(version.NewCommand())

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
