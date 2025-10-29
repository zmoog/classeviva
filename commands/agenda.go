package commands

import (
	"fmt"
	"sort"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/zmoog/classeviva/adapters/feedback"
	"github.com/zmoog/classeviva/adapters/spaggiari"
)

type ListAgendaCommand struct {
	Limit int
	Since time.Time
	Until time.Time
}

func (c ListAgendaCommand) ExecuteWith(uow UnitOfWork) error {
	entries, err := uow.Adapter.Agenda.List(c.Since, c.Until)
	if err != nil {
		return err
	}

	sort.Sort(AgendaEntriesByDate(entries))

	if c.Limit < len(entries) {
		entries = entries[:c.Limit]
	}

	return feedback.PrintResult(AgendaEntriesResult{Entries: entries})
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
	t.SetStyle(table.StyleColoredBright)
	t.SetColumnConfigs([]table.ColumnConfig{{Number: 1, AutoMerge: true}})

	t.AppendHeader(table.Row{"Time", "Subject", "Teacher", "Notes"})
	for _, e := range r.Entries {
		// Optimize the time formatting for visualization
		datetimeBegin, _ := time.Parse(time.RFC3339, e.DatetimeBegin)
		datetimeEnd, _ := time.Parse(time.RFC3339, e.DatetimeEnd)
		timeInt := fmt.Sprintf("%s - %s", datetimeBegin.Format("02-01-2006 15:04"), datetimeEnd.Format("02-01-2006 15:04"))
		if e.IsFullDay {
			timeInt = fmt.Sprintf("%s (all day)", datetimeBegin.Format("02-01-2006"))
		}

		// Use a placeholder for the subject if it's empty
		subject := "-"
		if e.Subject != "" {
			subject = e.Subject
		}

		t.AppendRow(table.Row{
			timeInt,
			subject,
			e.AuthorName,
			text.WrapSoft(e.Notes, 50),
		})
	}

	return t.Render()
}

// Data returns an interface holding the underlying data structure.
func (r AgendaEntriesResult) Data() interface{} {
	return r.Entries
}
