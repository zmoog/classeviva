package grades

import (
	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/commands"
)

func initSummarizeCommand() *cobra.Command {
	summarizeCommand := cobra.Command{
		Use:   "summarize",
		Short: "Summarize the grades on the portal",
		RunE:  runSummarizeCommand,
	}

	return &summarizeCommand
}

func runSummarizeCommand(cmd *cobra.Command, args []string) error {
	command := commands.SummarizeGradesCommand{}

	runner, err := commands.NewRunner()
	if err != nil {
		return err
	}

	err = runner.Run(command)
	if err != nil {
		return err
	}

	return nil
}
