package noticeboards_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zmoog/classeviva/adapters/config"
	"github.com/zmoog/classeviva/adapters/feedback"
	"github.com/zmoog/classeviva/adapters/spaggiari"
	"github.com/zmoog/classeviva/entrypoints/cli/cmd/noticeboards"
)

func TestReviewCommandDaysFlag(t *testing.T) {
	parent := noticeboards.NewCommand()
	cmd, _, err := parent.Find([]string{"review"})
	require.NoError(t, err)
	require.Equal(t, "review", cmd.Name())
	days, err := cmd.Flags().GetInt("days")
	require.NoError(t, err)
	assert.Equal(t, 7, days)
	require.NoError(t, cmd.ParseFlags([]string{"--days", "14"}))
	days, err = cmd.Flags().GetInt("days")
	require.NoError(t, err)
	assert.Equal(t, 14, days)

	// A new command tree has its own flag state; list gets no new flags.
	parent = noticeboards.NewCommand()
	cmd, _, err = parent.Find([]string{"review"})
	require.NoError(t, err)
	days, err = cmd.Flags().GetInt("days")
	require.NoError(t, err)
	assert.Equal(t, 7, days)
	list, _, err := parent.Find([]string{"list"})
	require.NoError(t, err)
	assert.Nil(t, list.Flags().Lookup("days"))
}

func TestReviewCommandRejectsInvalidInputBeforeCredentials(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want string
	}{
		{"negative days", []string{"review", "--days", "-1"}, "days must be zero or greater"},
		{"noninteger days", []string{"review", "--days", "tomorrow"}, "invalid argument"},
		{"positional argument", []string{"review", "unexpected"}, "unknown command"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			cmd := noticeboards.NewCommand()
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})
			cmd.SetArgs(tc.args)
			err := cmd.Execute()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.want)
			entries, err := os.ReadDir(home)
			require.NoError(t, err)
			assert.Empty(t, entries, "invalid input must not initialize profile configuration")
		})
	}
}

func TestReviewCommandUsesSelectedProfileAndOnlyListsNotices(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	cfg := config.DefaultConfig()
	cfg.AddProfile("synthetic-student", config.Profile{Username: "synthetic-user", Password: "synthetic-password"})
	require.NoError(t, config.Save(&cfg))
	now := time.Now().UTC()
	store := spaggiari.FilesystemLoaderStorer{Path: home, Profile: "synthetic-student"}
	require.NoError(t, store.Store(spaggiari.Identity{
		ID: "9000", Token: "synthetic-token",
		Release: now.Add(-time.Hour).Format(time.RFC3339),
		Expire:  now.Add(time.Hour).Format(time.RFC3339),
	}))

	requests := 0
	previous := http.DefaultTransport
	http.DefaultTransport = reviewTransport(func(req *http.Request) (*http.Response, error) {
		requests++
		if req.Method != http.MethodGet || req.URL.Path != "/rest/v1/students/9000/noticeboard" {
			return nil, fmt.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
		}
		assert.Equal(t, "synthetic-token", req.Header.Get("Z-Auth-Token"))
		body := `{"items":[
			{"pubId":41,"cntNum":"7","readStatus":true,"needSign":true},
			{"pubId":42,"readStatus":false}
		]}`
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString(body)), Header: http.Header{}}, nil
	})
	t.Cleanup(func() { http.DefaultTransport = previous })

	var stdout, stderr bytes.Buffer
	feedback.SetDefault(feedback.New(&stdout, &stderr, feedback.JSON))
	t.Cleanup(func() { feedback.SetDefault(feedback.Default()) })
	root := &cobra.Command{Use: "classeviva"}
	root.PersistentFlags().String("profile", "", "Profile")
	root.PersistentFlags().String("username", "", "Username")
	root.PersistentFlags().String("password", "", "Password")
	root.AddCommand(noticeboards.NewCommand())
	root.SetArgs([]string{"--profile", "synthetic-student", "noticeboards", "review", "--days", "14"})
	require.NoError(t, root.Execute())
	assert.Equal(t, 1, requests)
	var items []struct {
		Notice    spaggiari.Noticeboard `json:"notice"`
		Attention struct {
			Reasons []string `json:"reasons"`
		} `json:"attention"`
	}
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &items))
	require.Len(t, items, 1)
	assert.Equal(t, 41, items[0].Notice.ID)
	assert.True(t, items[0].Notice.Read)
	assert.Equal(t, []string{"signature_flag"}, items[0].Attention.Reasons)
	assert.Empty(t, stderr.String())
}

type reviewTransport func(*http.Request) (*http.Response, error)

func (transport reviewTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return transport(req)
}

func TestReviewCommandHelpIsReadOnly(t *testing.T) {
	var out bytes.Buffer
	cmd := noticeboards.NewCommand()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"review", "--help"})
	require.NoError(t, cmd.Execute())
	for _, want := range []string{"--days int", "(default 7)", "read-only", "Europe/Rome", "workflow"} {
		assert.Contains(t, out.String(), want)
	}
}
