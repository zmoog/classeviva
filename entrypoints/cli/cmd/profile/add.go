package profile

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/adapters/config"
)

var (
	addUsername string
	addPassword string
)

func initAddCommand() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Add a new profile",
		Args:  cobra.ExactArgs(1),
		RunE:  runAddCommand,
	}

	cmd.Flags().StringVarP(&addUsername, "username", "u", "", "Username for the profile")
	cmd.Flags().StringVar(&addPassword, "password", "", "Password for the profile")

	return cmd
}

func runAddCommand(cmd *cobra.Command, args []string) error {
	profileName := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if profile already exists
	if _, exists := cfg.GetProfile(profileName); exists {
		return fmt.Errorf("profile '%s' already exists", profileName)
	}

	// Get username and password (interactive if not provided via flags)
	username := addUsername
	password := addPassword

	reader := bufio.NewReader(os.Stdin)

	if username == "" {
		fmt.Print("Username: ")
		username, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read username: %w", err)
		}
		username = strings.TrimSpace(username)
	}

	if password == "" {
		fmt.Print("Password: ")
		password, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		password = strings.TrimSpace(password)
	}

	// Validate
	if username == "" || password == "" {
		return fmt.Errorf("username and password are required")
	}

	// Add profile
	cfg.AddProfile(profileName, config.Profile{
		Username: username,
		Password: password,
	})

	// If this is the first profile, make it default
	if len(cfg.Profiles) == 1 {
		cfg.DefaultProfile = profileName
		fmt.Printf("Profile '%s' added and set as default.\n", profileName)
	} else {
		fmt.Printf("Profile '%s' added successfully.\n", profileName)
	}

	// Save config
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}
