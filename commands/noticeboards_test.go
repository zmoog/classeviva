package commands_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zmoog/classeviva/adapters/feedback"
	"github.com/zmoog/classeviva/adapters/spaggiari"
	"github.com/zmoog/classeviva/commands"
)

func TestNoticeboardsJSONOutputPreservesResponseFields(t *testing.T) {
	data, err := os.ReadFile("../adapters/spaggiari/testdata/noticeboards.json")
	require.NoError(t, err)

	var response struct {
		Items []spaggiari.Noticeboard `json:"items"`
	}
	require.NoError(t, json.Unmarshal(data, &response))
	var expected struct {
		Items json.RawMessage `json:"items"`
	}
	require.NoError(t, json.Unmarshal(data, &expected))

	var stdout, stderr bytes.Buffer
	fb := feedback.New(&stdout, &stderr, feedback.JSON)
	require.NoError(t, fb.PrintResult(commands.NoticeboardsResult{Items: response.Items}))

	assert.JSONEq(t, string(expected.Items), stdout.String())
	assert.Empty(t, stderr.String())
}
