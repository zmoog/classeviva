package grades

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/adapters/spaggiari"
)

func initListCommand() *cobra.Command {
	summarizeCmd := cobra.Command{
		Use:   "list",
		Short: "List the grades on the portal",
		// Long:  "Summarize time entries in Toggl",
		RunE: runListCommand,
	}

	return &summarizeCmd
}

func runListCommand(cmd *cobra.Command, args []string) error {
	usernane := os.Getenv("CLASSEVIVA_USERNAME")
	password := os.Getenv("CLASSEVIVA_PASSWORD")

	adapter, err := spaggiari.From(usernane, password)
	if err != nil {
		fmt.Println(err)
	}

	grades, err := adapter.List()
	if err != nil {
		return err
	}

	sort.Sort(ByDate(grades))

	output, _ := json.MarshalIndent(grades[:3], "", "  ")
	fmt.Println(string(output))

	return nil
}

type ByDate []spaggiari.Grade

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Less(i, j int) bool { return a[i].Date > a[j].Date }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
