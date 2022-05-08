package version

import (
	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/commands"
)

func NewCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "version",
		Short: "Prints the application version",
		RunE:  runVersionCommand,
	}

	return &cmd
}

func runVersionCommand(cmd *cobra.Command, args []string) error {
	command := commands.VersionCommand{}

	runner, err := commands.NewRunner()
	if err != nil {
		return err
	}

	err = runner.Run(command)
	if err != nil {
		return err
	}

	return nil
}
