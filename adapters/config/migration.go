package config

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const (
	oldIdentityFile = "identity.json"
	defaultProfileName = "older-kid"
)

// MigrateIfNeeded checks for old identity.json and migrates to profile-based system
func MigrateIfNeeded() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	oldIdentityPath := filepath.Join(configDir, oldIdentityFile)
	newIdentityPath := filepath.Join(configDir, fmt.Sprintf("identity-%s.json", defaultProfileName))
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Check if migration is needed:
	// 1. Old identity.json exists
	// 2. Config file doesn't exist yet
	// 3. New profile-based identity doesn't exist
	oldExists := fileExists(oldIdentityPath)
	configExists := fileExists(configPath)
	newExists := fileExists(newIdentityPath)

	if !oldExists || configExists || newExists {
		// No migration needed
		return nil
	}

	log.Info("Migrating from single-student to multi-profile configuration...")

	// Rename identity.json to identity-older-kid.json
	if err := os.Rename(oldIdentityPath, newIdentityPath); err != nil {
		return fmt.Errorf("failed to rename identity file: %w", err)
	}
	log.Debugf("renamed %s to %s", oldIdentityPath, newIdentityPath)

	// Create a basic config with instructions
	// Note: We don't have the username/password from the old setup,
	// so we create an empty config with a message
	config := DefaultConfig()
	config.DefaultProfile = defaultProfileName

	// Check if we can get credentials from environment variables
	envUsername := os.Getenv("CLASSEVIVA_USERNAME")
	envPassword := os.Getenv("CLASSEVIVA_PASSWORD")

	if envUsername != "" && envPassword != "" {
		// Migrate env vars to default profile
		config.AddProfile(defaultProfileName, Profile{
			Username: envUsername,
			Password: envPassword,
		})
		log.Infof("Created profile '%s' with credentials from environment variables", defaultProfileName)
	}

	if err := Save(&config); err != nil {
		return fmt.Errorf("failed to save migrated config: %w", err)
	}

	// Print friendly migration message
	fmt.Println("\n" + migrationMessage())

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func migrationMessage() string {
	return `╔═══════════════════════════════════════════════════════════════════════════╗
║                     Classeviva Multi-Profile Setup                       ║
╚═══════════════════════════════════════════════════════════════════════════╝

Your configuration has been migrated to support multiple student profiles!

Your existing authentication token has been preserved as the 'older-kid' profile.

Next steps:
  1. Add your credentials to the config file:
     ~/.classeviva/config.yaml

  2. Add a profile for your second student:
     classeviva profile add younger-kid

  3. Switch between profiles:
     classeviva --profile older-kid grades list
     classeviva --profile younger-kid grades list

  4. Set a default profile (optional):
     classeviva profile set-default older-kid

For more information: classeviva profile --help
`
}
