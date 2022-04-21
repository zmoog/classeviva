package grades

import (
	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/commands"
)

var (
	limit int = 3
)

func initListCommand() *cobra.Command {
	listCommand := cobra.Command{
		Use:   "list",
		Short: "List the grades on the portal",
		// Long:  "Summarize time entries in Toggl",
		RunE: runListCommand,
	}

	listCommand.Flags().IntVarP(&limit, "limit", "l", limit, "Limit number of results")

	return &listCommand
}

func runListCommand(cmd *cobra.Command, args []string) error {
	runner, err := commands.NewRunner()
	if err != nil {
		return err
	}

	command := commands.ListGradesCommand{
		Limit: limit,
	}

	err = runner.Run(command)
	if err != nil {
		return err
	}

	return nil
}
