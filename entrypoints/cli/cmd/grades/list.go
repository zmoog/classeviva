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
		RunE:  runListCommand,
	}

	listCommand.Flags().IntVarP(&limit, "limit", "l", limit, "Limit number of results")

	return &listCommand
}

func runListCommand(cobraCmd *cobra.Command, args []string) error {
	command := commands.ListGradesCommand{
		Limit: limit,
	}

	// Get flags from parent command (persistent flags)
	profile, _ := cobraCmd.Flags().GetString("profile")
	username, _ := cobraCmd.Flags().GetString("username")
	password, _ := cobraCmd.Flags().GetString("password")

	runner, err := commands.NewRunner(commands.RunnerOptions{
		Username: username,
		Password: password,
		Profile:  profile,
	})
	if err != nil {
		return err
	}

	err = runner.Run(command)
	if err != nil {
		return err
	}

	return nil
}
