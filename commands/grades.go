package commands

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/zmoog/classeviva/adapters/spaggiari"
)

type ListGradesCommand struct {
	Limit int
}

func (c ListGradesCommand) Execute(uow UnitOfWork) error {

	grades, err := uow.Adapter.List()
	if err != nil {
		return err
	}

	sort.Sort(ByDate(grades))

	max := len(grades) - 1
	if c.Limit > 0 && c.Limit < max {
		max = c.Limit
	}

	output, _ := json.MarshalIndent(grades[:max], "", "  ")
	fmt.Println(string(output))

	return nil
}

type ByDate []spaggiari.Grade

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Less(i, j int) bool { return a[i].Date > a[j].Date }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
