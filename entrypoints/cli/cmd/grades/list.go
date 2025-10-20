package grades

import (
	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/commands"
	"github.com/zmoog/classeviva/entrypoints/cli/config"
)

var (
	limit int = 3
)

func initListCommand() *cobra.Command {
	listCommand := cobra.Command{
		Use:   "list",
		Short: "List the grades on the portal",
		RunE:  runListCommand,
	}

	listCommand.Flags().IntVarP(&limit, "limit", "l", limit, "Limit number of results")

	return &listCommand
}

func runListCommand(cobraCmd *cobra.Command, args []string) error {
	command := commands.ListGradesCommand{
		Limit: limit,
	}

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
