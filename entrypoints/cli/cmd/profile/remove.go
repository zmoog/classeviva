package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/adapters/config"
)

func initRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a profile",
		Args:  cobra.ExactArgs(1),
		RunE:  runRemoveCommand,
	}
}

func runRemoveCommand(cmd *cobra.Command, args []string) error {
	profileName := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Remove profile
	if removed := cfg.RemoveProfile(profileName); !removed {
		return fmt.Errorf("profile '%s' not found", profileName)
	}

	// Save config
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Profile '%s' removed successfully.\n", profileName)

	if cfg.DefaultProfile == "" && len(cfg.Profiles) > 0 {
		fmt.Println("\nTip: Set a default profile with:")
		fmt.Println("  classeviva profile set-default <name>")
	}

	return nil
}
