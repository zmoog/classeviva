package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/zmoog/classeviva/adapters/feedback"
)

type ReviewNoticeboardsCommand struct {
	Days int
}

func (c ReviewNoticeboardsCommand) ExecuteWith(uow UnitOfWork) error {
	if c.Days < 0 {
		return fmt.Errorf("days must be zero or greater")
	}
	notices, err := uow.Adapter.Noticeboards.List()
	if err != nil {
		return fmt.Errorf("failed to list noticeboards: %w", err)
	}
	items, err := ReviewNoticeboards(notices, time.Now(), c.Days)
	if err != nil {
		return err
	}
	return feedback.PrintResult(NoticeboardsReviewResult{Items: items})
}

type NoticeboardsReviewResult struct {
	Items []ReviewedNoticeboard
}

func (r NoticeboardsReviewResult) String() string {
	if len(r.Items) == 0 {
		return "No notices match the review criteria."
	}
	labels := map[string]string{
		"participation_flag":   "Participation flag set",
		"reply_flag":           "Reply flag set",
		"file_submission_flag": "File submission flag set",
		"signature_flag":       "Signature flag set",
		"unread":               "Unread",
		"expires_soon":         "expires soon",
	}
	t := table.NewWriter()
	t.AppendHeader(table.Row{"ID", "No.", "Read", "Valid until", "Attention", "Title"})
	for _, item := range r.Items {
		reasons := make([]string, 0, len(item.Attention.Reasons))
		for _, reason := range item.Attention.Reasons {
			reasons = append(reasons, labels[reason])
		}
		expiry := "Unknown"
		if days := item.Attention.DaysUntilExpiry; days != nil {
			switch {
			case *days < 0:
				expiry = item.Notice.ValidTo + " (expired)"
			case *days == 0:
				expiry = item.Notice.ValidTo + " (today)"
			case *days == 1:
				expiry = item.Notice.ValidTo + " (in 1 day)"
			default:
				expiry = fmt.Sprintf("%s (in %d days)", item.Notice.ValidTo, *days)
			}
		}
		t.AppendRow(table.Row{item.Notice.ID, item.Notice.Number, item.Notice.Read,
			expiry, strings.Join(reasons, "; "), item.Notice.Title})
	}
	return t.Render()
}

func (r NoticeboardsReviewResult) Data() interface{} {
	if r.Items == nil {
		return []ReviewedNoticeboard{}
	}
	return r.Items
}
