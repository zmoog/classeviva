package noticeboards

import (
	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/commands"
	"github.com/zmoog/classeviva/entrypoints/cli/config"
)

func initListCommand() *cobra.Command {
	listCommand := cobra.Command{
		Use:   "list",
		Short: "List the noticeboards on the portal",
		RunE:  runListCommand,
	}

	return &listCommand
}

func runListCommand(cobraCmd *cobra.Command, args []string) error {
	command := commands.ListNoticeboardsCommand{}

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
