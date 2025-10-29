package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/adapters/config"
)

func initSetDefaultCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set-default <name>",
		Short: "Set the default profile",
		Args:  cobra.ExactArgs(1),
		RunE:  runSetDefaultCommand,
	}
}

func runSetDefaultCommand(cmd *cobra.Command, args []string) error {
	profileName := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if profile exists
	if _, exists := cfg.GetProfile(profileName); !exists {
		return fmt.Errorf("profile '%s' not found", profileName)
	}

	// Set as default
	cfg.DefaultProfile = profileName

	// Save config
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Default profile set to '%s'\n", profileName)
	return nil
}
