package commands

import (
	"fmt"
	"sort"
	"strings"

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

func (r GradesResult) String() string {
	var sb strings.Builder
	for _, grade := range r.Grades {
		fmt.Fprintf(&sb, "%v %v %v", grade.Date, grade.Subject, grade.DisplaylValue)
		if grade.Notes != "" {
			fmt.Fprintf(&sb, " (%v)", grade.Notes)
		}
		fmt.Fprintf(&sb, "\n")
	}
	return sb.String()
}

func (r GradesResult) Data() interface{} {
	return r.Grades
}
