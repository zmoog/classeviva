package commands

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/zmoog/classeviva/adapters/spaggiari"
)

type ListNoticeboardsCommand struct{}

func (c ListNoticeboardsCommand) ExecuteWith(uow UnitOfWork) error {
	items, err := uow.Adapter.Noticeboards.List()
	if err != nil {
		return err
	}

	return uow.Feedback.PrintResult(NoticeboardsResult{Items: items})
}

type NoticeboardsResult struct {
	Items []spaggiari.Noticeboard
}

// String returns a string representation of the grades.
func (r NoticeboardsResult) String() string {
	if len(r.Items) == 0 {
		return "No noticeboards in this interval."
	}

	t := table.NewWriter()
	t.SetColumnConfigs([]table.ColumnConfig{{Number: 1, AutoMerge: true}})
	t.AppendHeader(table.Row{"PublicationDate", "Read", "Title"})

	for _, g := range r.Items {
		t.AppendRow(table.Row{g.PublicationDate, g.Read, g.Title})
	}

	return t.Render()
}

// Data returns an interface holding with a `[]spaggiari.Noticeboard` data structure.
func (r NoticeboardsResult) Data() interface{} {
	return r.Items
}
