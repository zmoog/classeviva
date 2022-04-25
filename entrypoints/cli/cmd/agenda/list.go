package agenda

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/commands"
)

var (
	limit int = 3
	since string
	until string
)

func initListCommand() *cobra.Command {
	listCommand := cobra.Command{
		Use:   "list",
		Short: "List the agenda entries",
		RunE:  runListCommand,
	}

	listCommand.Flags().IntVarP(&limit, "limit", "l", limit, "Limit number of results")
	listCommand.Flags().StringVarP(&since, "since", "s", time.Now().Format("2006-01-02"), "Day to summarize (format: YYYY-MM-DD)")
	listCommand.Flags().StringVarP(&until, "until", "u", time.Now().Add(3*24*time.Hour).Format("2006-01-02"), "Day to summarize (format: YYYY-MM-DD)")

	return &listCommand
}

func runListCommand(cmd *cobra.Command, args []string) error {
	_since, err := time.Parse("2006-01-02", since)
	if err != nil {
		return fmt.Errorf("invalid 'since' value: %w", err)
	}

	_until, err := time.Parse("2006-01-02", until)
	if err != nil {
		return fmt.Errorf("invalid 'until' value: %w", err)
	}

	command := commands.ListAgendaCommand{
		Limit: limit,
		Since: _since,
		Until: _until,
	}

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
