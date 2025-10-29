package commands_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zmoog/classeviva/commands"
)

type TestCommand struct{}

func (c TestCommand) ExecuteWith(uow commands.UnitOfWork) error {
	return nil
}

func TestRuner(t *testing.T) {

	t.Run("Fail without credentials", func(t *testing.T) {
		// Unset any existing environment variables to ensure a clean test environment
		t.Setenv("CLASSEVIVA_USERNAME", "")
		t.Setenv("CLASSEVIVA_PASSWORD", "")

		// Use a temporary directory for config to isolate test from user's actual config
		tempDir := t.TempDir()
		t.Setenv("HOME", tempDir)

		// With the new profile-based system, when no credentials are found via any method,
		// we expect an error indicating no credentials could be resolved
		_, err := commands.NewRunner(commands.RunnerOptions{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no credentials found")
	})

	t.Run("Execute command with UoW", func(t *testing.T) {
		// Use a temporary directory for config to isolate test from user's actual config
		tempDir := t.TempDir()
		t.Setenv("HOME", tempDir)

		t.Setenv("CLASSEVIVA_USERNAME", "test")
		t.Setenv("CLASSEVIVA_PASSWORD", "test")

		testCommand := TestCommand{}

		runner, err := commands.NewRunner(commands.RunnerOptions{})
		assert.Nil(t, err)

		err = runner.Run(testCommand)
		assert.Nil(t, err)
	})
}
