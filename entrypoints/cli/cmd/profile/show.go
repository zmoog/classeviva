package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/adapters/config"
)

func initShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show profile details (without password)",
		Args:  cobra.ExactArgs(1),
		RunE:  runShowCommand,
	}
}

func runShowCommand(cmd *cobra.Command, args []string) error {
	profileName := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get profile
	profile, exists := cfg.GetProfile(profileName)
	if !exists {
		return fmt.Errorf("profile '%s' not found", profileName)
	}

	// Show details
	fmt.Printf("Profile: %s\n", profileName)
	fmt.Printf("Username: %s\n", profile.Username)
	fmt.Printf("Password: %s\n", maskPassword(profile.Password))

	if profileName == cfg.DefaultProfile {
		fmt.Println("Default: Yes")
	} else {
		fmt.Println("Default: No")
	}

	return nil
}

func maskPassword(password string) string {
	if len(password) == 0 {
		return "(not set)"
	}
	return "********"
}
