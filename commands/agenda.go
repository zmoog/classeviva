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
}

func (c ListAgendaCommand) Execute(adapter spaggiari.Adapter) error {

	now := time.Now()
	since := now
	until := now.Add(5 * 24 * time.Hour)

	fmt.Println(&now, "---", now)
	fmt.Println(&since, "---", since)
	fmt.Println(&until, "---", until)

	entries, err := adapter.ListAgenda(since, until)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		return nil
	}

	sort.Sort(AgendaEntryByDate(entries))

	max := len(entries) - 1
	if c.Limit > 0 && c.Limit < max {
		max = c.Limit
	}

	output, _ := json.MarshalIndent(entries[:max], "", "  ")
	fmt.Println(string(output))

	return nil
}

type AgendaEntryByDate []spaggiari.AgendaEntry

func (a AgendaEntryByDate) Len() int           { return len(a) }
func (a AgendaEntryByDate) Less(i, j int) bool { return a[i].DatetimeBegin < a[j].DatetimeBegin }
func (a AgendaEntryByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
