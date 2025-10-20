package commands

import (
	"errors"
	"fmt"
	"os"

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

func NewRunner(cliUsername, cliPassword string) (Runner, error) {
	// Use CLI flags if provided, otherwise fall back to environment variables
	username := cliUsername
	password := cliPassword

	if username == "" {
		username = os.Getenv("CLASSEVIVA_USERNAME")
	}
	if password == "" {
		password = os.Getenv("CLASSEVIVA_PASSWORD")
	}

	if username == "" || password == "" {
		return Runner{}, errors.New("username and password must be provided via CLI flags (--username, --password) or environment variables (CLASSEVIVA_USERNAME, CLASSEVIVA_PASSWORD)")
	}

	identityStorePath, err := os.UserHomeDir()
	if err != nil {
		return Runner{}, fmt.Errorf("failed to get the user home dir: %w", err)
	}

	adapter, err := spaggiari.New(username, password, identityStorePath)
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
