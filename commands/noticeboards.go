package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/zmoog/classeviva/adapters/feedback"
	"github.com/zmoog/classeviva/adapters/spaggiari"
)

type ListNoticeboardsCommand struct{}

func (c ListNoticeboardsCommand) ExecuteWith(uow UnitOfWork) error {
	items, err := uow.Adapter.Noticeboards.List()
	if err != nil {
		return err
	}

	return feedback.PrintResult(NoticeboardsResult{Items: items})
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

type DownloadNoticeboardAttachmentCommand struct {
	PublicationID            int
	AttachmentSequenceNumber int
	OutputBasePath           string
}

func (c DownloadNoticeboardAttachmentCommand) ExecuteWith(uow UnitOfWork) error {
	// It seems there is no way to get a single noticeboard by ID,
	// so we have to list all of them and then find the one we are
	// looking for.
	items, err := uow.Adapter.Noticeboards.List()
	if err != nil {
		return fmt.Errorf("failed to list noticeboards: %w", err)
	}

	for _, item := range items {
		if item.ID == c.PublicationID {
			// Classeviva requires the noticeboard to be marked as read
			// before downloading the attachments.
			err := uow.Adapter.Noticeboards.SetAsRead(item.EventCode, item.ID)
			if err != nil {
				return fmt.Errorf("failed to set the noticeboard as read: %w", err)
			}

			var filenames []string

			for _, attachment := range item.Attachments {
				document, err := uow.Adapter.Noticeboards.DownloadAttachment(c.PublicationID, attachment.Number)
				if err != nil {
					return fmt.Errorf("failed to download the attachment: %w", err)
				}

				filename := path.Join(c.OutputBasePath, fmt.Sprintf("%d-%s", item.ID, attachment.Name))

				// write the document to the output file
				err = os.WriteFile(filename, document, 0644)
				if err != nil {
					return err
				}

				filenames = append(filenames, filename)
			}

			err = feedback.PrintResult(AttachmentDownloadResult{Filenames: filenames})
			if err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("noticeboard with ID %d not found", c.PublicationID)
}

type AttachmentDownloadResult struct {
	Filenames []string
}

func (r AttachmentDownloadResult) String() string {
	if len(r.Filenames) == 0 {
		return "No attachments to download."
	}

	t := table.NewWriter()
	t.SetColumnConfigs([]table.ColumnConfig{{Number: 1, AutoMerge: true}})
	t.AppendHeader(table.Row{"File"})

	for _, filename := range r.Filenames {
		t.AppendRow(table.Row{filename})
	}

	return t.Render()
}

func (r AttachmentDownloadResult) Data() interface{} {
	return r.Filenames
}
