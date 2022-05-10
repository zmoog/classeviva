package commands_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zmoog/classeviva/adapters/feedback"
	"github.com/zmoog/classeviva/commands"
	"github.com/zmoog/classeviva/mocks"
)

func TestVersion(t *testing.T) {
	t.Run("Text version", func(t *testing.T) {
		mockAdapter := mocks.Adapter{}

		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		fb := feedback.New(&stdout, &stderr, feedback.Text)
		feedback.SetDefault(fb)

		uow := commands.UnitOfWork{Adapter: &mockAdapter, Feedback: fb}

		cmd := commands.VersionCommand{}

		err := cmd.ExecuteWith(uow)
		assert.Nil(t, err)

		assert.Equal(t, "Classeviva CLI v0.0.0 (123) 2022-05-08 by zmoog", stdout.String())
		assert.Equal(t, "", stderr.String())
	})

	t.Run("JSON version", func(t *testing.T) {
		mockAdapter := mocks.Adapter{}

		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		fb := feedback.New(&stdout, &stderr, feedback.JSON)
		feedback.SetDefault(fb)

		uow := commands.UnitOfWork{Adapter: &mockAdapter, Feedback: fb}

		cmd := commands.VersionCommand{}

		err := cmd.ExecuteWith(uow)
		assert.Nil(t, err)

		expectedJSON := `{
  "version": "v0.0.0",
  "commit": "123",
  "date": "2022-05-08",
  "built_by": "zmoog"
}`

		assert.Equal(t, expectedJSON, stdout.String())
		assert.Equal(t, "", stderr.String())
	})
}
