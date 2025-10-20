package commands_test

import (
	"errors"
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
		t.Setenv("CLASSEVIVA_USERNAME", "")
		t.Setenv("CLASSEVIVA_PASSWORD", "")

		expected := errors.New("username and password must be provided via CLI flags (--username, --password) or environment variables (CLASSEVIVA_USERNAME, CLASSEVIVA_PASSWORD)")
		_, err := commands.NewRunner("", "")
		assert.Equal(t, expected, err)
	})

	t.Run("Execute command with environment variables", func(t *testing.T) {
		t.Setenv("CLASSEVIVA_USERNAME", "test")
		t.Setenv("CLASSEVIVA_PASSWORD", "test")

		testCommand := TestCommand{}

		runner, err := commands.NewRunner("", "")
		assert.Nil(t, err)

		err = runner.Run(testCommand)
		assert.Nil(t, err)
	})

	t.Run("Execute command with CLI flags", func(t *testing.T) {
		testCommand := TestCommand{}

		runner, err := commands.NewRunner("testuser", "testpass")
		assert.Nil(t, err)

		err = runner.Run(testCommand)
		assert.Nil(t, err)
	})

	t.Run("CLI flags take precedence over environment variables", func(t *testing.T) {
		t.Setenv("CLASSEVIVA_USERNAME", "envuser")
		t.Setenv("CLASSEVIVA_PASSWORD", "envpass")

		testCommand := TestCommand{}

		runner, err := commands.NewRunner("cliuser", "clipass")
		assert.Nil(t, err)

		err = runner.Run(testCommand)
		assert.Nil(t, err)
	})
}
