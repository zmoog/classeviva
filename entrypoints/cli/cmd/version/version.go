package version

import (
	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/commands"
	"github.com/zmoog/classeviva/entrypoints/cli/config"
)

func NewCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "version",
		Short: "Prints the application version",
		RunE:  runVersionCommand,
	}

	return &cmd
}

func runVersionCommand(cobraCmd *cobra.Command, args []string) error {
	command := commands.VersionCommand{}

	username, password := config.GetCredentials()
	runner, err := commands.NewRunner(username, password)
	if err != nil {
		return err
	}

	err = runner.Run(command)
	if err != nil {
		return err
	}

	return nil
}
