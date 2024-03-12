package noticeboards

import (
	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/commands"
)

func initListCommand() *cobra.Command {
	listCommand := cobra.Command{
		Use:   "list",
		Short: "List the noticeboards on the portal",
		RunE:  runListCommand,
	}

	return &listCommand
}

func runListCommand(cmd *cobra.Command, args []string) error {
	command := commands.ListNoticeboardsCommand{}

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
