package commands

import (
	"fmt"
	"os"

	"github.com/zmoog/classeviva/adapters/config"
	"github.com/zmoog/classeviva/adapters/spaggiari"
)

type Runner struct {
	uow UnitOfWork
}

func (r Runner) Run(command Command) error {
	err := command.ExecuteWith(r.uow)
	if err != nil {
		return err
	}

	return nil
}

// RunnerOptions contains options for creating a new Runner
type RunnerOptions struct {
	Username string
	Password string
	Profile  string
}

func NewRunner(opts RunnerOptions) (Runner, error) {
	// Run migration check on first use
	if err := config.MigrateIfNeeded(); err != nil {
		return Runner{}, fmt.Errorf("failed to migrate config: %w", err)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return Runner{}, fmt.Errorf("failed to load config: %w", err)
	}

	// Resolve credentials using the priority chain
	creds, err := config.ResolveCredentials(config.ResolverOptions{
		Username: opts.Username,
		Password: opts.Password,
		Profile:  opts.Profile,
		Config:   cfg,
	})
	if err != nil {
		return Runner{}, fmt.Errorf("failed to resolve credentials: %w", err)
	}

	// Validate credentials
	if err := config.ValidateCredentials(creds); err != nil {
		return Runner{}, fmt.Errorf("invalid credentials: %w", err)
	}

	identityStorePath, err := os.UserHomeDir()
	if err != nil {
		return Runner{}, fmt.Errorf("failed to get the user home dir: %w", err)
	}

	adapter, err := spaggiari.New(creds.Username, creds.Password, identityStorePath, creds.Profile)
	if err != nil {
		return Runner{}, err
	}

	runner := Runner{
		uow: UnitOfWork{
			Adapter: adapter,
			// Feedback: feedback.Default(),
		},
	}

	return runner, nil
}
