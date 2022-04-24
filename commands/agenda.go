package commands

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/zmoog/classeviva/adapters/spaggiari"
)

type ListAgendaCommand struct {
	Limit int
	Since time.Time
	Until time.Time
}

func (c ListAgendaCommand) Execute(uow UnitOfWork) error {
	entries, err := uow.Adapter.ListAgenda(c.Since, c.Until)
	if err != nil {
		return err
	}

	sort.Sort(AgendaEntriesByDate(entries))

	if c.Limit < len(entries) {
		entries = entries[:c.Limit]
	}

	output, _ := json.MarshalIndent(entries, "", "  ")
	uow.Feedback.Println(string(output))

	return nil
}

type AgendaEntriesByDate []spaggiari.AgendaEntry

func (a AgendaEntriesByDate) Len() int           { return len(a) }
func (a AgendaEntriesByDate) Less(i, j int) bool { return a[i].DatetimeBegin < a[j].DatetimeBegin }
func (a AgendaEntriesByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
