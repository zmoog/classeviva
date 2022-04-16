package grades

import (
	"fmt"

	"github.com/spf13/cobra"
)

func initListCommand() *cobra.Command {
	summarizeCmd := cobra.Command{
		Use:   "list",
		Short: "List the grades on the portal",
		// Long:  "Summarize time entries in Toggl",
		RunE: runListCommand,
	}

	// summarizeCmd.Flags().StringVarP(&dayDate, "day", "d", "", "Day to summarize (format: YYYY-MM-DD)")

	// err := summarizeCmd.MarkFlagRequired("day")
	// if err != nil {
	// panic(err)
	// }

	return &summarizeCmd
}

func runListCommand(cmd *cobra.Command, args []string) error {
	// mb := messagebus.ForCLI()
	// const dateFormat = "2006-01-02"

	// from := "2021-08-20T00:00:00.000Z"
	// to := "2021-08-20T23:59:59.999Z"
	// t, err := time.Parse(time.RFC3339, str)

	// from, err := time.Parse(dateFormat, fromDate)
	// if err != nil {
	// 	return err
	// }
	// to, err := time.Parse(dateFormat, toDate)
	// if err != nil {
	// 	return err
	// }
	// day, err := time.Parse(dateFormat, dayDate)
	// if err != nil {
	// 	return err
	// }

	// ctx := context.Background()
	// command := commands.SummarizeEntriesCommand{
	// 	From: day,
	// 	To:   day,
	// }

	// return mb.Process(ctx, command)

	fmt.Println("faking the list of grades")

	return nil
}
