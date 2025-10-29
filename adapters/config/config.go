package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	log "github.com/sirupsen/logrus"
)

// Config represents the application configuration
type Config struct {
	Profiles       map[string]Profile `koanf:"profiles"`
	DefaultProfile string             `koanf:"default_profile"`
}

// Profile represents a student profile with credentials
type Profile struct {
	Username string `koanf:"username"`
	Password string `koanf:"password"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Profiles:       make(map[string]Profile),
		DefaultProfile: "",
	}
}

// Load loads the configuration from file and environment variables
func Load() (*Config, error) {
	k := koanf.New(".")

	// 1. Load defaults
	if err := k.Load(structs.Provider(DefaultConfig(), "koanf"), nil); err != nil {
		return nil, fmt.Errorf("failed to load default config: %w", err)
	}

	// 2. Load config file if it exists
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	if _, err := os.Stat(configPath); err == nil {
		log.Debugf("loading config from: %s", configPath)
		if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	} else {
		log.Debugf("config file not found at: %s", configPath)
	}

	// 3. Load environment variables (CLASSEVIVA_PROFILES_<profile>_USERNAME, etc.)
	if err := k.Load(env.Provider("CLASSEVIVA_", ".", func(s string) string {
		// Transform CLASSEVIVA_PROFILES_OLDER_KID_USERNAME -> profiles.older-kid.username
		s = strings.Replace(strings.ToLower(s), "classeviva_", "", 1)
		s = strings.Replace(s, "_", ".", -1)
		return s
	}), nil); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	// Unmarshal into Config struct
	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// Save saves the configuration to the config file
func Save(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Ensure directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create koanf instance and load the config struct
	k := koanf.New(".")
	if err := k.Load(structs.Provider(config, "koanf"), nil); err != nil {
		return fmt.Errorf("failed to load config into koanf: %w", err)
	}

	// Marshal to YAML
	yamlBytes, err := k.Marshal(yaml.Parser())
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %w", err)
	}

	// Write to file with restrictive permissions
	if err := os.WriteFile(configPath, yamlBytes, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	log.Debugf("config saved to: %s", configPath)
	return nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(home, ".classeviva", "config.yaml"), nil
}

// GetConfigDir returns the path to the config directory
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(home, ".classeviva"), nil
}

// GetProfile retrieves a profile by name from the config
func (c *Config) GetProfile(name string) (Profile, bool) {
	profile, exists := c.Profiles[name]
	return profile, exists
}

// AddProfile adds or updates a profile in the config
func (c *Config) AddProfile(name string, profile Profile) {
	if c.Profiles == nil {
		c.Profiles = make(map[string]Profile)
	}
	c.Profiles[name] = profile
}

// RemoveProfile removes a profile from the config
func (c *Config) RemoveProfile(name string) bool {
	if _, exists := c.Profiles[name]; exists {
		delete(c.Profiles, name)
		// Clear default if it was the removed profile
		if c.DefaultProfile == name {
			c.DefaultProfile = ""
		}
		return true
	}
	return false
}

// ListProfiles returns a list of all profile names
func (c *Config) ListProfiles() []string {
	names := make([]string, 0, len(c.Profiles))
	for name := range c.Profiles {
		names = append(names, name)
	}
	return names
}
