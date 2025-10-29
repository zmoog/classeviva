package noticeboards

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/commands"
)

var (
	publicationID int
	outputDir     string
)

func initDownloadCommand() *cobra.Command {
	downloadCommand := cobra.Command{
		Use:   "download",
		Short: "Download the noticeboard attachment from the portal",
		RunE:  runDownloadCommand,
	}

	downloadCommand.Flags().IntVar(&publicationID, "publication_id", publicationID, "Publication ID to download the attachment from")
	downloadCommand.Flags().StringVarP(&outputDir, "output-filename", "o", outputDir, "Output directory for the attachment(s)")

	return &downloadCommand
}

func runDownloadCommand(cobraCmd *cobra.Command, args []string) error {
	if publicationID == 0 {
		return fmt.Errorf("publication_id is required")
	}

	if outputDir == "" {
		outputDir = "."
	}

	command := commands.DownloadNoticeboardAttachmentCommand{
		PublicationID:  publicationID,
		OutputBasePath: outputDir,
	}

	// Get flags from parent command (persistent flags)
	profile, _ := cobraCmd.Flags().GetString("profile")
	username, _ := cobraCmd.Flags().GetString("username")
	password, _ := cobraCmd.Flags().GetString("password")

	runner, err := commands.NewRunner(commands.RunnerOptions{
		Username: username,
		Password: password,
		Profile:  profile,
	})
	if err != nil {
		return err
	}

	err = runner.Run(command)
	if err != nil {
		return err
	}

	return nil
}
