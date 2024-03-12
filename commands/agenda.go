package commands

import (
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/zmoog/classeviva/adapters/spaggiari"
)

type ListAgendaCommand struct {
	Limit int
	Since time.Time
	Until time.Time
}

func (c ListAgendaCommand) ExecuteWith(uow UnitOfWork) error {
	// entries, err := uow.Adapter.ListAgenda(c.Since, c.Until)
	// if err != nil {
	// 	return err
	// }

	// sort.Sort(AgendaEntriesByDate(entries))

	// if c.Limit < len(entries) {
	// 	entries = entries[:c.Limit]
	// }

	// return feedback.PrintResult(AgendaEntriesResult{Entries: entries})
	return nil
}

type AgendaEntriesByDate []spaggiari.AgendaEntry

func (a AgendaEntriesByDate) Len() int           { return len(a) }
func (a AgendaEntriesByDate) Less(i, j int) bool { return a[i].DatetimeBegin < a[j].DatetimeBegin }
func (a AgendaEntriesByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type AgendaEntriesResult struct {
	Entries []spaggiari.AgendaEntry
}

// String returns a string representation of the agenda entries.
func (r AgendaEntriesResult) String() string {
	if len(r.Entries) == 0 {
		return "No entries in this interval."
	}

	t := table.NewWriter()
	t.SetColumnConfigs([]table.ColumnConfig{{Number: 1, AutoMerge: true}})

	t.AppendHeader(table.Row{"Begin", "End", "Subject", "Teacher", "Notes"})
	for _, e := range r.Entries {
		t.AppendRow(table.Row{
			e.DatetimeBegin,
			e.DatetimeEnd,
			e.Subject,
			e.AuthorName,
			e.Notes,
		})
	}

	return t.Render()
}

// Data returns an interface holding the underlying data structure.
func (r AgendaEntriesResult) Data() interface{} {
	return r.Entries
}
