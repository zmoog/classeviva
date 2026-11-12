package spaggiari_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zmoog/classeviva/adapters/spaggiari"
)

func TestNoticeboardPreservesResponseFields(t *testing.T) {
	// All data is synthetic. True workflow/change flags exercise preservation;
	// their meaning and positive production examples have not been verified.
	data, err := os.ReadFile("testdata/noticeboards.json")
	require.NoError(t, err)

	var envelope struct {
		Items []spaggiari.Noticeboard `json:"items"`
	}
	require.NoError(t, json.Unmarshal(data, &envelope))
	require.Len(t, envelope.Items, 2)

	notice := envelope.Items[0]
	assert.Equal(t, 501, notice.ContentID)
	assert.Equal(t, "007", notice.Number)
	assert.Equal(t, "2026-11-11", notice.ValidFrom)
	assert.Equal(t, "2027-02-28", notice.ValidTo)
	assert.True(t, notice.HasChanged)
	assert.True(t, notice.NeedJoin)
	assert.True(t, notice.NeedReply)
	assert.True(t, notice.NeedFile)
	assert.True(t, notice.NeedSign)
	assert.Equal(t, "501", notice.EventID)
	assert.Equal(t, "2026-11-11 13:43:59", notice.InsertedAt)

	encoded, err := json.Marshal(envelope)
	require.NoError(t, err)
	assert.JSONEq(t, string(data), string(encoded))
}
