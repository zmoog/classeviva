package grades

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "grades",
		Short: "Grade registerd",
	}

	cmd.AddCommand(initListCommand())

	return &cmd
}
