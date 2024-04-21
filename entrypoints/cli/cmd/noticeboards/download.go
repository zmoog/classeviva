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

	downloadCommand.Flags().IntVarP(&publicationID, "publication_id", "p", publicationID, "Publication ID to download the attachment from")
	// downloadCommand.Flags().IntVarP(&attachmentSequenceNumber, "sequence", "s", attachmentSequenceNumber, "Attachment sequence number")
	downloadCommand.Flags().StringVarP(&outputDir, "output-filename", "o", outputDir, "Output directory for the attachment(s)")

	return &downloadCommand
}

func runDownloadCommand(cmd *cobra.Command, args []string) error {
	if publicationID == 0 {
		return fmt.Errorf("pubblication_id is required")
	}
	// if attachmentSequenceNumber == 0 {
	// 	return fmt.Errorf("sequence is required")
	// }

	if outputDir == "" {
		outputDir = "."
	}

	command := commands.DownloadNoticeboardAttachmentCommand{
		PublicationID: publicationID,
		// AttachmentSequenceNumber: attachmentSequenceNumber,
		OutputBasePath: outputDir,
	}

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
