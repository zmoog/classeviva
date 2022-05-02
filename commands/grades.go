package commands

import (
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/zmoog/classeviva/adapters/feedback"
	"github.com/zmoog/classeviva/adapters/spaggiari"
)

type ListGradesCommand struct {
	Limit int
}

func (c ListGradesCommand) ExecuteWith(uow UnitOfWork) error {

	grades, err := uow.Adapter.List()
	if err != nil {
		return err
	}

	sort.Sort(ByDate(grades))

	if c.Limit < len(grades) {
		grades = grades[:c.Limit]
	}

	return feedback.PrintResult(GradesResult{Grades: grades})
}

type ByDate []spaggiari.Grade

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Less(i, j int) bool { return a[i].Date > a[j].Date }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type GradesResult struct {
	Grades []spaggiari.Grade
}

// String returns a string representation of the grades.
func (r GradesResult) String() string {
	if len(r.Grades) == 0 {
		return "No grades in this interval."
	}

	t := table.NewWriter()
	t.SetColumnConfigs([]table.ColumnConfig{{Number: 1, AutoMerge: true}})
	t.AppendHeader(table.Row{"Date", "Grade", "Subject", "Notes"})

	for _, g := range r.Grades {
		value := g.DisplaylValue
		switch g.Color {
		case "green":
			value = text.FgGreen.Sprint(value)
		case "red":
			value = text.FgRed.Sprint(value)
		case "blue":
			value = text.FgBlue.Sprint(value)
		}
		t.AppendRow(table.Row{g.Date, value, g.Subject, g.Notes})
	}

	return t.Render()
}

// Data returns an interface holding with a `[]spaggiari.Grade` data structure.
func (r GradesResult) Data() interface{} {
	return r.Grades
}
