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

	t.Run("Fail without environment variables", func(t *testing.T) {
		t.Setenv("CLASSEVIVA_USERNAME", "")
		t.Setenv("CLASSEVIVA_PASSWORD", "")

		expected := errors.New("CLASSEVIVA_USERNAME or CLASSEVIVA_PASSWORD environment variables are empty")
		_, err := commands.NewRunner()
		assert.Equal(t, expected, err)
	})

	t.Run("Execute command with UoW", func(t *testing.T) {
		t.Setenv("CLASSEVIVA_USERNAME", "test")
		t.Setenv("CLASSEVIVA_PASSWORD", "test")

		testCommand := TestCommand{}
		// mockCommand.On("ExecuteWith", mock.AnythingOfType("commands.UnitOfWork")).Return(nil)

		runner, err := commands.NewRunner()
		assert.Nil(t, err)

		err = runner.Run(testCommand)
		assert.Nil(t, err)
	})
}
