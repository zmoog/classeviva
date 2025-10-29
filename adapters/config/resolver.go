package config

import (
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// Credentials represents resolved user credentials
type Credentials struct {
	Username string
	Password string
	Profile  string // The profile name used (for identity caching)
}

// ResolverOptions contains options for credential resolution
type ResolverOptions struct {
	// CLI flags (highest priority)
	Username string
	Password string
	Profile  string

	// Config (loaded from file/env)
	Config *Config
}

// ResolveCredentials resolves credentials using the priority chain:
// 1. CLI flags (--username, --password)
// 2. Profile (--profile or default_profile from config)
// 3. Environment variables (CLASSEVIVA_USERNAME, CLASSEVIVA_PASSWORD)
// 4. Error if none found
func ResolveCredentials(opts ResolverOptions) (Credentials, error) {
	var creds Credentials

	// Priority 1: CLI flags for username/password (highest priority)
	if opts.Username != "" && opts.Password != "" {
		log.Debug("using credentials from CLI flags")
		creds.Username = opts.Username
		creds.Password = opts.Password
		creds.Profile = "cli-override" // Special profile name for CLI overrides
		return creds, nil
	}

	// Priority 2: Profile-based credentials
	if opts.Config != nil {
		profileName := opts.Profile
		if profileName == "" {
			profileName = opts.Config.DefaultProfile
		}

		if profileName != "" {
			profile, exists := opts.Config.GetProfile(profileName)
			if !exists {
				return Credentials{}, fmt.Errorf("profile '%s' not found in config", profileName)
			}

			if profile.Username != "" && profile.Password != "" {
				log.Debugf("using credentials from profile: %s", profileName)
				creds.Username = profile.Username
				creds.Password = profile.Password
				creds.Profile = profileName
				return creds, nil
			}
		}
	}

	// Priority 3: Environment variables (backward compatibility)
	envUsername := os.Getenv("CLASSEVIVA_USERNAME")
	envPassword := os.Getenv("CLASSEVIVA_PASSWORD")
	if envUsername != "" && envPassword != "" {
		log.Debug("using credentials from environment variables")
		creds.Username = envUsername
		creds.Password = envPassword
		creds.Profile = "env-override" // Special profile name for env var overrides
		return creds, nil
	}

	// No credentials found
	return Credentials{}, errors.New("no credentials found: provide via --username/--password flags, config profile, or CLASSEVIVA_USERNAME/CLASSEVIVA_PASSWORD environment variables")
}

// ValidateCredentials checks if credentials are complete
func ValidateCredentials(creds Credentials) error {
	if creds.Username == "" {
		return errors.New("username is empty")
	}
	if creds.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}
