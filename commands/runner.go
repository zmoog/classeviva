package commands

import (
	"errors"
	"os"

	"github.com/zmoog/classeviva/adapters/spaggiari"
)

type Runner struct {
	adapter spaggiari.Adapter
}

func (r Runner) Run(command Command) error {
	err := command.Execute(r.adapter)
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

	adapter, err := spaggiari.From(usernane, password)
	if err != nil {
		return Runner{}, err
	}

	runner := Runner{
		adapter: adapter,
	}

	return runner, nil
}
