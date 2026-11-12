package noticeboards

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/commands"
)

func initReviewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "review",
		Short: "Show notices worth reviewing",
		Long: `Show notices with any workflow flag set, or unread notices whose publication
validity expires within the selected calendar-day window (including today).

This is a read-only heuristic view: it does not mark notices read, download
attachments, or complete actions. Workflow flags are not verified pending tasks,
and publication expiry is not necessarily an action deadline.

Flagged notices remain visible even when read, expired, or missing an expiry date.
Dates are compared in Europe/Rome. Results show the reasons for inclusion.`,
		Args: cobra.NoArgs,
		RunE: runReviewCommand,
	}
	cmd.Flags().Int("days", 7, "Calendar days ahead for unread expiry signals (0 means today only)")
	return cmd
}

func runReviewCommand(cmd *cobra.Command, args []string) error {
	days, err := cmd.Flags().GetInt("days")
	if err != nil {
		return err
	}
	if days < 0 {
		return fmt.Errorf("days must be zero or greater")
	}
	profile, _ := cmd.Flags().GetString("profile")
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	runner, err := commands.NewRunner(commands.RunnerOptions{
		Username: username,
		Password: password,
		Profile:  profile,
	})
	if err != nil {
		return err
	}
	return runner.Run(commands.ReviewNoticeboardsCommand{Days: days})
}
