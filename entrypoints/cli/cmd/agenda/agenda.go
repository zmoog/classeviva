package agenda

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "agenda",
		Short: "Agenda",
	}

	cmd.AddCommand(initListCommand())

	return &cmd
}
