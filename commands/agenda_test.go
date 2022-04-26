package commands_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zmoog/classeviva/adapters/feedback"
	"github.com/zmoog/classeviva/adapters/spaggiari"
	"github.com/zmoog/classeviva/commands"
	"github.com/zmoog/classeviva/mocks"
)

func TestListAgendaCommand(t *testing.T) {
	t.Run("Empty agenda items list", func(t *testing.T) {
		mockAdapter := mocks.Adapter{}
		mockAdapter.On(
			"ListAgenda",
			mock.AnythingOfType("time.Time"),
			mock.AnythingOfType("time.Time"),
		).Return([]spaggiari.AgendaEntry{}, nil)

		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		fb := feedback.New(&stdout, &stderr, feedback.Text)
		feedback.SetDefault(fb)

		uow := commands.UnitOfWork{Adapter: &mockAdapter, Feedback: fb}

		cmd := commands.ListAgendaCommand{}

		err := cmd.ExecuteWith(uow)
		assert.Nil(t, err)
		assert.Equal(t, stdout.String(), "[]\n")
		assert.Equal(t, stderr.String(), "")

		mockAdapter.AssertExpectations(t)
	})
}
