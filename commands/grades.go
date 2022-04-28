package commands

import (
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
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

	max := len(grades) - 1
	if c.Limit > 0 && c.Limit < max {
		max = c.Limit
	}

	return feedback.PrintResult(GradesResult{Grades: grades[:max]})
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
	t := table.NewWriter()

	t.SetColumnConfigs([]table.ColumnConfig{{Number: 1, AutoMerge: true}})
	// t.Style().Options.SeparateRows = true
	t.AppendHeader(table.Row{"Date", "Grade", "Subject", "Notes"})

	for _, g := range r.Grades {
		t.AppendRow(table.Row{g.Date, g.DisplaylValue, g.Subject, g.Notes})
	}

	return t.Render()
}

// Data returns an interface holding with a `[]spaggiari.Grade` data structure.
func (r GradesResult) Data() interface{} {
	return r.Grades
}
