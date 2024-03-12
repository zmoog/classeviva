package noticeboards

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "noticeboards",
		Short: "Noticeboards",
	}

	cmd.AddCommand(initListCommand())

	return &cmd
}
