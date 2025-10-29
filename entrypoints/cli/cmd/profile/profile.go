package profile

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage student profiles",
		Long:  "Manage multiple student profiles for accessing Classeviva",
	}

	profileCmd.AddCommand(initListCommand())
	profileCmd.AddCommand(initAddCommand())
	profileCmd.AddCommand(initSetDefaultCommand())
	profileCmd.AddCommand(initRemoveCommand())
	profileCmd.AddCommand(initShowCommand())

	return profileCmd
}
