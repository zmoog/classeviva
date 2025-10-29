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

func runListCommand(cobraCmd *cobra.Command, args []string) error {
	command := commands.ListNoticeboardsCommand{}

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
