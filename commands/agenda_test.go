package commands_test

import (
	"bytes"
	"os"
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

		cmd := commands.ListAgendaCommand{Limit: 10}

		err := cmd.ExecuteWith(uow)
		assert.Nil(t, err)
		assert.Equal(t, stdout.String(), "No entries in this interval.")
		assert.Equal(t, stderr.String(), "")

		mockAdapter.AssertExpectations(t)
	})

	t.Run("List 5 agenda entries", func(t *testing.T) {
		entries := []spaggiari.AgendaEntry{}
		if err := UnmarshalFrom("testdata/agenda.json", &entries); err != nil {
			t.Error(err)
		}

		mockAdapter := mocks.Adapter{}
		mockAdapter.On(
			"ListAgenda",
			mock.AnythingOfType("time.Time"),
			mock.AnythingOfType("time.Time"),
		).Return(entries, nil)

		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		fb := feedback.New(&stdout, &stderr, feedback.Text)
		feedback.SetDefault(fb)

		uow := commands.UnitOfWork{Adapter: &mockAdapter, Feedback: fb}
		cmd := commands.ListAgendaCommand{Limit: 10}

		err := cmd.ExecuteWith(uow)
		assert.Nil(t, err)

		expected, err := os.ReadFile("testdata/agenda.out.txt")
		if err != nil {
			t.Errorf("can't read test data from %v: %v", "testdata/agenda.out.txt", err)
		}

		assert.Equal(t, stdout.String(), string(expected))
		assert.Equal(t, stderr.String(), "")

		mockAdapter.AssertExpectations(t)
	})
}
