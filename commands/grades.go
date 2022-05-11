package commands

import (
	"fmt"
	"sort"
	"strings"

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

type SummarizeGradesCommand struct{}

func (c SummarizeGradesCommand) ExecuteWith(uow UnitOfWork) error {
	grades, err := uow.Adapter.List()
	if err != nil {
		return err
	}

	sort.Sort(ByDateAsc(grades))

	summary := map[string][]spaggiari.Grade{}

	for _, grade := range grades {
		if _, exists := summary[grade.Subject]; !exists {
			summary[grade.Subject] = []spaggiari.Grade{}
		}
		summary[grade.Subject] = append(summary[grade.Subject], grade)
	}

	return feedback.PrintResult(GradeSummaryResult{Summary: summary})
}

type GradeSummaryResult struct {
	Summary map[string][]spaggiari.Grade
}

func (r GradeSummaryResult) String() string {
	t := table.NewWriter()
	t.SetColumnConfigs([]table.ColumnConfig{{Number: 1, AutoMerge: true}})
	t.AppendHeader(table.Row{"Subject", "Grades", "Avg", "Trend"})

	for subject, grades := range r.Summary {
		gradesDisplay := []string{}
		sum := 0.0
		average := 0.0
		trend := "="

		for _, grade := range grades {
			// blue grades do not count for
			// grade average
			if grade.Color == "blue" {
				continue
			}

			sum += grade.DecimalValue

			gradesDisplay = append(gradesDisplay, grade.DisplaylValue)
			newAverage := sum / float64(len(gradesDisplay))

			switch {
			case newAverage > average:
				trend = text.FgGreen.Sprint("+")
				// trend = "+"
			case newAverage < average:
				trend = text.FgRed.Sprint("-")
				// trend = "-"
			default:
				trend = "="
			}

			// fmt.Println(sum, average, newAverage)
			average = newAverage
		}

		t.AppendRow(table.Row{
			subject,
			strings.Join(gradesDisplay, ", "),
			fmt.Sprintf("%.2f", average),
			// average,
			trend,
		})
	}

	t.SortBy([]table.SortBy{{Name: "Avg", Mode: table.DscNumeric}})
	// t.SetColumnConfigs([]table.ColumnConfig{{}})

	return t.Render()
}

func (r GradeSummaryResult) Data() interface{} {
	return r
}

type ByDateAsc []spaggiari.Grade

func (a ByDateAsc) Len() int           { return len(a) }
func (a ByDateAsc) Less(i, j int) bool { return a[i].Date < a[j].Date }
func (a ByDateAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByDate []spaggiari.Grade

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Less(i, j int) bool { return a[i].Date > a[j].Date }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
