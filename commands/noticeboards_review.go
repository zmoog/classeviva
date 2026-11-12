package commands

import (
	"fmt"
	"sort"
	"time"

	"github.com/zmoog/classeviva/adapters/spaggiari"
)

// ReviewedNoticeboard keeps the API data separate from derived review signals.
type ReviewedNoticeboard struct {
	Notice    spaggiari.Noticeboard `json:"notice"`
	Attention NoticeboardAttention  `json:"attention"`
	priority  int
}

type NoticeboardAttention struct {
	Reasons         []string `json:"reasons"`
	DaysUntilExpiry *int     `json:"days_until_expiry"`
}

// ReviewNoticeboards selects workflow-flagged notices and unread notices expiring
// within days (inclusive). Flagged notices are included regardless of read state
// or expiry. Dates are compared as Europe/Rome calendar dates, not durations.
// The input is not modified; a nil error always returns a non-nil slice.
func ReviewNoticeboards(notices []spaggiari.Noticeboard, now time.Time, days int) ([]ReviewedNoticeboard, error) {
	if days < 0 {
		return nil, fmt.Errorf("days must be zero or greater")
	}
	location, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		return nil, fmt.Errorf("load school timezone: %w", err)
	}
	local := now.In(location)
	// UTC midnights represent calendar dates only; this avoids 23/25-hour DST days.
	today := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, time.UTC)
	items := make([]ReviewedNoticeboard, 0)
	for _, notice := range notices {
		attention := NoticeboardAttention{Reasons: make([]string, 0)}
		if expiry, err := time.Parse(time.DateOnly, notice.ValidTo); err == nil {
			remaining := int(expiry.Unix()/86400 - today.Unix()/86400)
			attention.DaysUntilExpiry = &remaining
		}
		for _, flag := range []struct {
			set    bool
			reason string
		}{
			{notice.NeedJoin, "participation_flag"},
			{notice.NeedReply, "reply_flag"},
			{notice.NeedFile, "file_submission_flag"},
			{notice.NeedSign, "signature_flag"},
		} {
			if flag.set {
				attention.Reasons = append(attention.Reasons, flag.reason)
			}
		}
		flagged := len(attention.Reasons) > 0
		expiringUnread := !notice.Read && attention.DaysUntilExpiry != nil &&
			*attention.DaysUntilExpiry >= 0 && *attention.DaysUntilExpiry <= days
		if !flagged && !expiringUnread {
			continue
		}
		if !notice.Read {
			attention.Reasons = append(attention.Reasons, "unread")
		}
		if expiringUnread {
			attention.Reasons = append(attention.Reasons, "expires_soon")
		}
		priority := 2
		if flagged {
			priority = 1
			if expiringUnread {
				priority = 0
			}
		}
		items = append(items, ReviewedNoticeboard{Notice: notice, Attention: attention, priority: priority})
	}
	// Within a priority group, earliest expiry first, unknown expiry last.
	// Preserve server order for ties rather than assuming an API sort contract.
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].priority != items[j].priority {
			return items[i].priority < items[j].priority
		}
		a, b := items[i].Attention.DaysUntilExpiry, items[j].Attention.DaysUntilExpiry
		if a == nil {
			return false
		}
		return b == nil || *a < *b
	})
	return items, nil
}
