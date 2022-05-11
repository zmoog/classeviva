package commands_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zmoog/classeviva/adapters/feedback"
	"github.com/zmoog/classeviva/adapters/spaggiari"
	"github.com/zmoog/classeviva/commands"
	"github.com/zmoog/classeviva/mocks"
)

func TestListGradesCommand(t *testing.T) {
	t.Run("Empty grades list", func(t *testing.T) {
		mockAdapter := mocks.Adapter{}
		mockAdapter.On(
			"List",
		).Return([]spaggiari.Grade{}, nil)

		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		fb := feedback.New(&stdout, &stderr, feedback.Text)
		feedback.SetDefault(fb)

		uow := commands.UnitOfWork{Adapter: &mockAdapter, Feedback: fb}

		cmd := commands.ListGradesCommand{Limit: 100}

		err := cmd.ExecuteWith(uow)
		assert.Nil(t, err)
		assert.Equal(t, stdout.String(), "No grades in this interval.")
		assert.Equal(t, stderr.String(), "")

		mockAdapter.AssertExpectations(t)
	})

	t.Run("List 1 grade", func(t *testing.T) {
		entries := []spaggiari.Grade{}
		if err := UnmarshalFrom("testdata/grades.json", &entries); err != nil {
			t.Error(err)
		}

		mockAdapter := mocks.Adapter{}
		mockAdapter.On("List").Return(entries, nil)

		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		fb := feedback.New(&stdout, &stderr, feedback.Text)
		feedback.SetDefault(fb)

		uow := commands.UnitOfWork{Adapter: &mockAdapter, Feedback: fb}
		cmd := commands.ListGradesCommand{Limit: 10}

		err := cmd.ExecuteWith(uow)
		assert.Nil(t, err)

		expected, err := os.ReadFile("testdata/grades.out.txt")
		if err != nil {
			t.Errorf("can't read test data from %v: %v", "testdata/grades.out.txt", err)
		}

		assert.Equal(t, stdout.String(), string(expected))
		assert.Equal(t, stderr.String(), "")

		mockAdapter.AssertExpectations(t)
	})
}

func TestSummarize(t *testing.T) {
	t.Run("Summarize grades", func(t *testing.T) {
		entries := []spaggiari.Grade{}
		if err := UnmarshalFrom("testdata/grades-summarize.json", &entries); err != nil {
			t.Error(err)
		}

		mockAdapter := mocks.Adapter{}
		mockAdapter.On("List").Return(entries, nil)

		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		fb := feedback.New(&stdout, &stderr, feedback.Text)
		feedback.SetDefault(fb)

		uow := commands.UnitOfWork{Adapter: &mockAdapter, Feedback: fb}
		cmd := commands.SummarizeGradesCommand{}

		err := cmd.ExecuteWith(uow)
		assert.Nil(t, err)

		expected, err := ioutil.ReadFile("testdata/grades-summarize.out.txt")
		if err != nil {
			t.Errorf("can't read test data from %v: %v", "testdata/grades-summarize.out.txtt", err)
		}

		assert.Equal(t, string(expected), stdout.String())
		assert.Equal(t, "", stderr.String())

		// should I test for this method to be called?
		// mockAdapter.AssertExpectations(t)
	})
}
