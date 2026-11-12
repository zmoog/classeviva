package commands_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zmoog/classeviva/adapters/spaggiari"
	"github.com/zmoog/classeviva/commands"
)

func TestReviewNoticeboardsSelectsIndependentAttentionSignals(t *testing.T) {
	now := time.Date(2026, time.November, 12, 12, 0, 0, 0, time.UTC)
	notices := []spaggiari.Noticeboard{
		{ID: 1, Read: false, ValidTo: "2026-11-14"},
		{ID: 2, Read: true, NeedSign: true, ValidTo: "2027-06-30"},
		{ID: 3, Read: false, NeedReply: true, ValidTo: "2026-11-15"},
		{ID: 4, Read: false, ValidTo: "2026-12-31"},
		{ID: 5, Read: true, ValidTo: "2026-11-14"},
		{ID: 6, Read: false, ValidTo: "2026-11-11"},
	}

	items, err := commands.ReviewNoticeboards(notices, now, 14)
	require.NoError(t, err)
	require.Len(t, items, 3)
	// Combined signals first, then workflow flags, then unread expiry alone.
	assert.Equal(t, 3, items[0].Notice.ID)
	assert.Equal(t, []string{"reply_flag", "unread", "expires_soon"}, items[0].Attention.Reasons)
	require.NotNil(t, items[0].Attention.DaysUntilExpiry)
	assert.Equal(t, 3, *items[0].Attention.DaysUntilExpiry)
	assert.Equal(t, 2, items[1].Notice.ID)
	assert.Equal(t, []string{"signature_flag"}, items[1].Attention.Reasons)
	assert.Equal(t, 1, items[2].Notice.ID)
	assert.Equal(t, []string{"unread", "expires_soon"}, items[2].Attention.Reasons)
	// Assessment must not mark notices read or reorder the caller's slice.
	assert.Equal(t, 1, notices[0].ID)
	assert.False(t, notices[0].Read)
}

func TestReviewNoticeboardsExpiryAndWorkflowCases(t *testing.T) {
	cases := []struct {
		name      string
		expiry    string
		flag      string
		read      bool
		days      int
		want      bool
		known     bool
		remaining int
	}{
		{name: "today is included with zero days", expiry: "2026-11-12", days: 0, want: true, known: true},
		{name: "window end is inclusive", expiry: "2026-11-26", days: 14, want: true, known: true, remaining: 14},
		{name: "outside window", expiry: "2026-11-27", days: 14},
		{name: "expired unread is not due soon", expiry: "2026-11-11", days: 14},
		{name: "read without flags is excluded", expiry: "2026-11-13", read: true, days: 14},
		{name: "missing expiry without flags", days: 14},
		{name: "invalid expiry without flags", expiry: "not-a-date", days: 14},
		{name: "timestamp is not a date", expiry: "2026-11-13T00:00:00Z", days: 14},
		{name: "join alone", flag: "join", read: true, want: true},
		{name: "reply alone", flag: "reply", read: true, want: true},
		{name: "file alone", flag: "file", read: true, want: true},
		{name: "sign alone", flag: "sign", read: true, want: true},
		{name: "expired flagged notice is retained", expiry: "2026-11-11", flag: "sign", read: true, want: true, known: true, remaining: -1},
		{name: "flagged notice ignores window", expiry: "2026-12-31", flag: "reply", want: true, known: true, remaining: 49},
		{name: "flag survives invalid expiry", expiry: "2026-02-30", flag: "file", want: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			notice := spaggiari.Noticeboard{
				Read: tc.read, ValidTo: tc.expiry,
				NeedJoin: tc.flag == "join", NeedReply: tc.flag == "reply",
				NeedFile: tc.flag == "file", NeedSign: tc.flag == "sign",
			}
			items, err := commands.ReviewNoticeboards([]spaggiari.Noticeboard{notice},
				time.Date(2026, time.November, 12, 12, 0, 0, 0, time.UTC), tc.days)
			require.NoError(t, err)
			if !tc.want {
				assert.Empty(t, items)
				return
			}
			require.Len(t, items, 1)
			assert.Equal(t, notice, items[0].Notice)
			if tc.known {
				require.NotNil(t, items[0].Attention.DaysUntilExpiry)
				assert.Equal(t, tc.remaining, *items[0].Attention.DaysUntilExpiry)
			} else {
				assert.Nil(t, items[0].Attention.DaysUntilExpiry)
			}
		})
	}
}

func TestReviewNoticeboardsUsesRomeCalendarDates(t *testing.T) {
	cases := []struct {
		name, now, expiry string
		days              int
	}{
		{"winter midnight", "2026-11-12T23:30:00Z", "2026-11-13", 0},
		{"summer midnight", "2026-08-01T22:30:00Z", "2026-08-02", 0},
		{"spring DST transition", "2026-03-28T23:30:00Z", "2026-03-30", 1},
		{"autumn DST transition", "2026-10-24T22:30:00Z", "2026-10-26", 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			now, err := time.Parse(time.RFC3339, tc.now)
			require.NoError(t, err)
			items, err := commands.ReviewNoticeboards([]spaggiari.Noticeboard{{ValidTo: tc.expiry}}, now, tc.days)
			require.NoError(t, err)
			require.Len(t, items, 1)
			require.NotNil(t, items[0].Attention.DaysUntilExpiry)
			assert.Equal(t, tc.days, *items[0].Attention.DaysUntilExpiry)
		})
	}
}

func TestReviewNoticeboardsSortsExpiryWithinPriority(t *testing.T) {
	notices := []spaggiari.Noticeboard{
		{ID: 1, NeedSign: true},
		{ID: 2, Read: true, NeedFile: true, ValidTo: "2026-11-20"},
		{ID: 3, Read: true, NeedReply: true, ValidTo: "2026-11-19"},
		{ID: 4, Read: true, NeedSign: true, ValidTo: "2026-11-19"},
		{ID: 5, Read: true, NeedJoin: true, ValidTo: "invalid"},
	}
	items, err := commands.ReviewNoticeboards(notices, time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC), 7)
	require.NoError(t, err)
	ids := make([]int, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.Notice.ID)
	}
	assert.Equal(t, []int{3, 4, 2, 1, 5}, ids)
}

func TestReviewNoticeboardsEmptyAndInvalidWindow(t *testing.T) {
	items, err := commands.ReviewNoticeboards(nil, time.Now(), 7)
	require.NoError(t, err)
	assert.NotNil(t, items)
	assert.Empty(t, items)
	_, err = commands.ReviewNoticeboards(nil, time.Now(), -1)
	assert.EqualError(t, err, "days must be zero or greater")
}
