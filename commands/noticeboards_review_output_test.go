package commands_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zmoog/classeviva/adapters/feedback"
	"github.com/zmoog/classeviva/adapters/spaggiari"
	"github.com/zmoog/classeviva/commands"
)

func TestNoticeboardsReviewTextExplainsSignals(t *testing.T) {
	notices := []spaggiari.Noticeboard{
		{ID: 1001, Number: "007", Title: "Synthetic expired workflow", Read: true,
			NeedJoin: true, NeedReply: true, NeedFile: true, NeedSign: true, ValidTo: "2026-11-11"},
		{ID: 1002, Title: "Synthetic unread notice", ValidTo: "2026-11-13"},
		{ID: 1003, Title: "Synthetic undated workflow", Read: true, NeedSign: true},
		{ID: 1004, Title: "Synthetic expires today", ValidTo: "2026-11-12"},
		{ID: 1005, Title: "Synthetic expires in two days", ValidTo: "2026-11-14"},
	}
	items, err := commands.ReviewNoticeboards(notices, time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC), 7)
	require.NoError(t, err)
	text := (commands.NoticeboardsReviewResult{Items: items}).String()
	for _, want := range []string{
		"1001", "007", "READ", "VALID UNTIL", "ATTENTION", "TITLE",
		"2026-11-11 (expired)", "2026-11-13 (in 1 day)", "Unknown",
		"2026-11-12 (today)", "2026-11-14 (in 2 days)",
		"Participation flag set", "Reply flag set", "File submission flag set", "Signature flag set",
		"Unread", "expires soon", "Synthetic expired workflow",
	} {
		assert.Contains(t, text, want)
	}
	assert.NotContains(t, text, "pending")
	assert.NotContains(t, text, "overdue")
}

func TestNoticeboardsReviewJSONPreservesNoticeAndAddsReasons(t *testing.T) {
	data, err := os.ReadFile("../adapters/spaggiari/testdata/noticeboards.json")
	require.NoError(t, err)
	var response struct {
		Items []spaggiari.Noticeboard `json:"items"`
	}
	var raw struct {
		Items []json.RawMessage `json:"items"`
	}
	require.NoError(t, json.Unmarshal(data, &response))
	require.NoError(t, json.Unmarshal(data, &raw))
	items, err := commands.ReviewNoticeboards(response.Items, time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC), 14)
	require.NoError(t, err)
	var stdout, stderr bytes.Buffer
	fb := feedback.New(&stdout, &stderr, feedback.JSON)
	require.NoError(t, fb.PrintResult(commands.NoticeboardsReviewResult{Items: items}))
	assert.JSONEq(t, fmt.Sprintf(`[{"notice": %s, "attention": {
		"reasons": ["participation_flag", "reply_flag", "file_submission_flag", "signature_flag"],
		"days_until_expiry": 108
	}}]`, raw.Items[0]), stdout.String())
	assert.Empty(t, stderr.String())
}

func TestNoticeboardsReviewEmptyOutput(t *testing.T) {
	result := commands.NoticeboardsReviewResult{}
	assert.Equal(t, "No notices match the review criteria.", result.String())
	data, err := json.Marshal(result.Data())
	require.NoError(t, err)
	assert.JSONEq(t, "[]", string(data))
}

// Embedding the interface leaves all operations except List unimplemented:
// any accidental read-marking or download call fails the read-only command test.
type reviewNoticeboardsReceiver struct {
	spaggiari.NoticeboardsReceiver
	items []spaggiari.Noticeboard
	err   error
}

func (r reviewNoticeboardsReceiver) List() ([]spaggiari.Noticeboard, error) {
	return r.items, r.err
}

func TestReviewNoticeboardsCommandIsReadOnly(t *testing.T) {
	var stdout, stderr bytes.Buffer
	feedback.SetDefault(feedback.New(&stdout, &stderr, feedback.JSON))
	t.Cleanup(func() { feedback.SetDefault(feedback.Default()) })
	adapter := spaggiari.Adapter{Noticeboards: reviewNoticeboardsReceiver{items: []spaggiari.Noticeboard{
		{ID: 1, Read: true, NeedSign: true},
		{ID: 2, Read: false},
	}}}
	err := (commands.ReviewNoticeboardsCommand{Days: 14}).ExecuteWith(commands.UnitOfWork{Adapter: adapter})
	require.NoError(t, err)
	var items []struct {
		Notice    spaggiari.Noticeboard         `json:"notice"`
		Attention commands.NoticeboardAttention `json:"attention"`
	}
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &items))
	require.Len(t, items, 1)
	assert.Equal(t, 1, items[0].Notice.ID)
	assert.True(t, items[0].Notice.Read)
	assert.Equal(t, []string{"signature_flag"}, items[0].Attention.Reasons)
	assert.Nil(t, items[0].Attention.DaysUntilExpiry)
	assert.Contains(t, stdout.String(), `"days_until_expiry": null`)
	assert.Empty(t, stderr.String())
}

func TestReviewNoticeboardsCommandErrors(t *testing.T) {
	want := errors.New("synthetic API failure")
	adapter := spaggiari.Adapter{Noticeboards: reviewNoticeboardsReceiver{err: want}}
	err := (commands.ReviewNoticeboardsCommand{Days: 7}).ExecuteWith(commands.UnitOfWork{Adapter: adapter})
	assert.ErrorIs(t, err, want)
	// Invalid options are rejected before any adapter access.
	err = (commands.ReviewNoticeboardsCommand{Days: -1}).ExecuteWith(commands.UnitOfWork{})
	assert.EqualError(t, err, "days must be zero or greater")
}
