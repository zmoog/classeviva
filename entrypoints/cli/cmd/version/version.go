package version

import (
	"fmt"

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

	// Version command doesn't need credentials, pass empty options
	runner, err := commands.NewRunner(commands.RunnerOptions{})
	if err != nil {
		// Version command doesn't require auth, so just print version without runner
		// This allows "classeviva version" to work without configuration
		result := commands.VersionResult{
			Version: "v0.0.0",
			Commit:  "123",
			Date:    "2022-05-08",
			BuiltBy: "zmoog",
		}
		fmt.Println(result.String())
		return nil
	}

	err = runner.Run(command)
	if err != nil {
		return err
	}

	return nil
}
