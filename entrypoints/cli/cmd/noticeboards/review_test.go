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

func TestReviewCommandFlags(t *testing.T) {
	parent := noticeboards.NewCommand()
	cmd, _, err := parent.Find([]string{"review"})
	require.NoError(t, err)
	require.Equal(t, "review", cmd.Name())
	days, err := cmd.Flags().GetInt("days")
	require.NoError(t, err)
	assert.Equal(t, 7, days)
	unread, err := cmd.Flags().GetBool("unread")
	require.NoError(t, err)
	assert.False(t, unread)
	require.NoError(t, cmd.ParseFlags([]string{"--days", "14", "--unread"}))
	days, err = cmd.Flags().GetInt("days")
	require.NoError(t, err)
	assert.Equal(t, 14, days)
	unread, err = cmd.Flags().GetBool("unread")
	require.NoError(t, err)
	assert.True(t, unread)

	// A new command tree has its own flag state; list gets no new flags.
	parent = noticeboards.NewCommand()
	cmd, _, err = parent.Find([]string{"review"})
	require.NoError(t, err)
	days, err = cmd.Flags().GetInt("days")
	require.NoError(t, err)
	assert.Equal(t, 7, days)
	unread, err = cmd.Flags().GetBool("unread")
	require.NoError(t, err)
	assert.False(t, unread)
	list, _, err := parent.Find([]string{"list"})
	require.NoError(t, err)
	assert.Nil(t, list.Flags().Lookup("days"))
	assert.Nil(t, list.Flags().Lookup("unread"))
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
			{"pubId":42,"readStatus":false},
			{"pubId":43,"readStatus":false,"needReply":true}
		]}`
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString(body)), Header: http.Header{}}, nil
	})
	t.Cleanup(func() { http.DefaultTransport = previous })

	for _, tc := range []struct {
		name    string
		flags   []string
		wantIDs []int
	}{
		{"default includes read flags", nil, []int{41, 43}},
		{"unread excludes read flags", []string{"--unread"}, []int{43}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			requests = 0
			var stdout, stderr bytes.Buffer
			feedback.SetDefault(feedback.New(&stdout, &stderr, feedback.JSON))
			t.Cleanup(func() { feedback.SetDefault(feedback.Default()) })
			root := &cobra.Command{Use: "classeviva"}
			root.PersistentFlags().String("profile", "", "Profile")
			root.PersistentFlags().String("username", "", "Username")
			root.PersistentFlags().String("password", "", "Password")
			root.AddCommand(noticeboards.NewCommand())
			args := []string{"--profile", "synthetic-student", "noticeboards", "review", "--days", "14"}
			root.SetArgs(append(args, tc.flags...))
			require.NoError(t, root.Execute())
			assert.Equal(t, 1, requests)
			var items []struct {
				Notice    spaggiari.Noticeboard `json:"notice"`
				Attention struct {
					Reasons []string `json:"reasons"`
				} `json:"attention"`
			}
			require.NoError(t, json.Unmarshal(stdout.Bytes(), &items))
			wantReasons := map[int][]string{41: {"signature_flag"}, 43: {"reply_flag", "unread"}}
			ids := make([]int, 0, len(items))
			for _, item := range items {
				ids = append(ids, item.Notice.ID)
				assert.Equal(t, item.Notice.ID == 41, item.Notice.Read)
				assert.Equal(t, wantReasons[item.Notice.ID], item.Attention.Reasons)
			}
			assert.Equal(t, tc.wantIDs, ids)
			assert.Empty(t, stderr.String())
		})
	}
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
	for _, want := range []string{"--days int", "--unread", "(default 7)", "read-only", "Europe/Rome", "workflow"} {
		assert.Contains(t, out.String(), want)
	}
}
