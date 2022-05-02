package feedback_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zmoog/classeviva/adapters/feedback"
	"github.com/zmoog/classeviva/adapters/spaggiari"
	"github.com/zmoog/classeviva/commands"
)

func TestPrint(t *testing.T) {

	t.Run("Println message to standard output", func(t *testing.T) {
		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		fb := feedback.New(&stdout, &stderr, feedback.Text)
		feedback.SetDefault(fb)

		feedback.Println("welcome to this f* world!")

		assert.Equal(t, "welcome to this f* world!\n", stdout.String())
		assert.Equal(t, "", stderr.String())
	})
}
func TestError(t *testing.T) {

	t.Run("Print message to standard error", func(t *testing.T) {
		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		fb := feedback.New(&stdout, &stderr, feedback.Text)
		feedback.SetDefault(fb)

		feedback.Error("Doh!")

		assert.Equal(t, "", stdout.String())
		assert.Equal(t, "Doh!\n", stderr.String())
	})
}

func TestPrintResult(t *testing.T) {

	t.Run("Print result in plain text", func(t *testing.T) {
		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		fb := feedback.New(&stdout, &stderr, feedback.Text)
		feedback.SetDefault(fb)

		grades := []spaggiari.Grade{
			{
				Subject:       "COMPORTAMENTO",
				Date:          "2022-04-22",
				DecimalValue:  7.25,
				DisplaylValue: "7+",
				Color:         "green",
				Notes:         "comportamento della settimana",
			},
			{
				Subject:       "SCIENZE",
				Date:          "2022-04-22",
				DecimalValue:  7,
				DisplaylValue: "7",
				Color:         "green",
				Description:   " orale",
			},
		}

		err := feedback.PrintResult(commands.GradesResult{
			Grades: grades,
		})
		assert.NoError(t, err)

		assert.Equal(t,
			"+------------+-------+---------------+-------------------------------+\n"+
				"| DATE       | GRADE | SUBJECT       | NOTES                         |\n"+
				"+------------+-------+---------------+-------------------------------+\n"+
				"| 2022-04-22 | \x1b[32m7+\x1b[0m    | COMPORTAMENTO | comportamento della settimana |\n"+
				"|            | \x1b[32m7\x1b[0m     | SCIENZE       |                               |\n"+
				"+------------+-------+---------------+-------------------------------+",
			stdout.String(),
		)
		assert.Equal(t, "", stderr.String())
	})

	t.Run("Print result in JSON", func(t *testing.T) {
		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		fb := feedback.New(&stdout, &stderr, feedback.JSON)
		feedback.SetDefault(fb)

		grades := []spaggiari.Grade{
			{
				Subject:       "COMPORTAMENTO",
				Date:          "2022-04-22",
				DecimalValue:  7.25,
				DisplaylValue: "7+",
				Color:         "green",
				Notes:         "comportamento della settimana",
			},
			{
				Subject:       "SCIENZE",
				Date:          "2022-04-22",
				DecimalValue:  7,
				DisplaylValue: "7",
				Color:         "green",
				Description:   " orale",
			},
		}

		err := feedback.PrintResult(commands.GradesResult{
			Grades: grades,
		})
		assert.NoError(t, err)

		assert.Equal(t, `[
  {
    "subjectDesc": "COMPORTAMENTO",
    "evtDate": "2022-04-22",
    "decimalValue": 7.25,
    "displayValue": "7+",
    "color": "green",
    "notesForFamily": "comportamento della settimana"
  },
  {
    "subjectDesc": "SCIENZE",
    "evtDate": "2022-04-22",
    "decimalValue": 7,
    "displayValue": "7",
    "color": "green",
    "skillValueDesc": " orale"
  }
]`, stdout.String())
		assert.Equal(t, "", stderr.String())
	})
}
