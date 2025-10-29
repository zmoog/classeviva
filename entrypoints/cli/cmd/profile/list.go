package profile

import (
	"fmt"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/zmoog/classeviva/adapters/config"
)

func initListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all profiles",
		RunE:  runListCommand,
	}
}

func runListCommand(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.Profiles) == 0 {
		fmt.Println("No profiles configured.")
		fmt.Println("\nAdd a profile:")
		fmt.Println("  classeviva profile add <name>")
		return nil
	}

	// Get profiles and sort by name
	names := cfg.ListProfiles()
	sort.Strings(names)

	t := table.NewWriter()
	t.AppendHeader(table.Row{"Profile", "Username", "Default"})

	for _, name := range names {
		profile, _ := cfg.GetProfile(name)
		isDefault := ""
		if name == cfg.DefaultProfile {
			isDefault = "*"
		}
		t.AppendRow(table.Row{name, profile.Username, isDefault})
	}

	fmt.Println(t.Render())
	return nil
}
