package commands

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/zmoog/classeviva/adapters/spaggiari"
)

type ListAgendaCommand struct {
	Limit int
	Since time.Time
	Until time.Time
}

func (c ListAgendaCommand) Execute(adapter spaggiari.Adapter) error {
	entries, err := adapter.ListAgenda(c.Since, c.Until)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		fmt.Println("No agenda entries for the given since/until interval")
		return nil
	}

	sort.Sort(AgendaEntriesByDate(entries))

	max := len(entries) - 1
	if c.Limit > 0 && c.Limit < max {
		max = c.Limit
	}

	output, _ := json.MarshalIndent(entries[:max], "", "  ")
	fmt.Println(string(output))

	return nil
}

type AgendaEntriesByDate []spaggiari.AgendaEntry

func (a AgendaEntriesByDate) Len() int           { return len(a) }
func (a AgendaEntriesByDate) Less(i, j int) bool { return a[i].DatetimeBegin < a[j].DatetimeBegin }
func (a AgendaEntriesByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
