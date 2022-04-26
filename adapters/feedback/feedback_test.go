package feedback_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zmoog/classeviva/adapters/feedback"
)

func TestPrint(t *testing.T) {

	t.Run("Println message to standard output", func(t *testing.T) {
		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		fb := feedback.New(&stdout, &stderr)

		fb.Println("welcome to this f* world!")

		assert.Equal(t, "welcome to this f* world!\n", stdout.String())
		assert.Equal(t, "", stderr.String())
	})
}
func TestError(t *testing.T) {

	t.Run("Print message to standard error", func(t *testing.T) {
		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		fb := feedback.New(&stdout, &stderr)

		fb.Error("Doh!")

		assert.Equal(t, "", stdout.String())
		assert.Equal(t, "Doh!\n", stderr.String())
	})
}
