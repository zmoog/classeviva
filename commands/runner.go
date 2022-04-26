package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/zmoog/classeviva/adapters/feedback"
	"github.com/zmoog/classeviva/adapters/spaggiari"
)

type Runner struct {
	uow UnitOfWork
}

func (r Runner) Run(command Command) error {
	err := command.Execute(r.uow)
	if err != nil {
		return err
	}

	return nil
}

func NewRunner() (Runner, error) {
	usernane := os.Getenv("CLASSEVIVA_USERNAME")
	password := os.Getenv("CLASSEVIVA_PASSWORD")
	if usernane == "" || password == "" {
		return Runner{}, errors.New("CLASSEVIVA_USERNAME or CLASSEVIVA_PASSWORD environment variables are empty")
	}

	identityStorePath, err := os.UserHomeDir()
	if err != nil {
		return Runner{}, fmt.Errorf("failed to get the user home dir: %w", err)
	}

	adapter, err := spaggiari.From(usernane, password, identityStorePath)
	if err != nil {
		return Runner{}, err
	}

	runner := Runner{
		uow: UnitOfWork{
			Adapter:  adapter,
			Feedback: feedback.Default(),
		},
	}

	return runner, nil
}
