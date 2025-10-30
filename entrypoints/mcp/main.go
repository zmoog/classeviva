package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
	"github.com/zmoog/classeviva/adapters/config"
	"github.com/zmoog/classeviva/adapters/spaggiari"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home dir: %w", err)
	}

	cfg, err := config.Load()
	if err != nil {
		configPath, _ := config.GetConfigPath()
		return fmt.Errorf("failed to load configuration: %w\n\nPlease create a configuration file at: %s\n\nExample format:\n%s",
			err, configPath, getExampleConfig())
	}

	if len(cfg.Profiles) == 0 {
		configPath, _ := config.GetConfigPath()
		return fmt.Errorf("no profiles configured\n\nPlease add profiles to: %s\n\nExample format:\n%s",
			configPath, getExampleConfig())
	}

	profileAdapters := make(map[string]spaggiari.Adapter)
	for profileName, profile := range cfg.Profiles {
		adapter, err := spaggiari.New(profile.Username, profile.Password, homeDir, profileName)
		if err != nil {
			return fmt.Errorf("failed to create adapter for profile %s: %w", profileName, err)
		}
		profileAdapters[profileName] = adapter
	}

	s := server.NewMCPServer(
		"Classeviva MCP Server",
		"1.0.0",
		server.WithLogging(),
	)

	s.AddTool(mcp.NewTool("list_profiles",
		mcp.WithDescription("List all configured student profiles"),
		mcp.WithString("format",
			mcp.Description("Output format: 'text' or 'json'"),
			mcp.DefaultString("json"),
		),
	), listProfilesHandler(cfg))

	s.AddTool(mcp.NewTool("list_grades",
		mcp.WithDescription("List grades for a specific student profile"),
		mcp.WithString("profile",
			mcp.Description("The profile name from list_profiles"),
			mcp.Required(),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of grades to return (default: all)"),
		),
	), listGradesHandler(cfg, profileAdapters))

	s.AddTool(mcp.NewTool("list_agenda",
		mcp.WithDescription("List agenda items (homework, events) for a specific student profile"),
		mcp.WithString("profile",
			mcp.Description("The profile name from list_profiles"),
			mcp.Required(),
		),
		mcp.WithString("since",
			mcp.Description("Start date in YYYY-MM-DD format (default: today)"),
		),
		mcp.WithString("until",
			mcp.Description("End date in YYYY-MM-DD format (default: 30 days from now)"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of items to return (default: all)"),
		),
	), listAgendaHandler(cfg, profileAdapters))

	s.AddTool(mcp.NewTool("list_noticeboards",
		mcp.WithDescription("List noticeboard items (announcements, circulars) for a specific student profile"),
		mcp.WithString("profile",
			mcp.Description("The profile name from list_profiles"),
			mcp.Required(),
		),
	), listNoticeboardsHandler(cfg, profileAdapters))

	if err := server.ServeStdio(s); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

func listProfilesHandler(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		format := request.GetString("format", "json")

		profiles := make([]map[string]string, 0, len(cfg.Profiles))
		for profileName := range cfg.Profiles {
			profiles = append(profiles, map[string]string{
				"name": profileName,
			})
		}

		var output string
		if format == "json" {
			data, err := json.MarshalIndent(profiles, "", "  ")
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal profiles: %v", err)), nil
			}
			output = string(data)
		} else {
			output = "Available profiles:\n"
			for profileName := range cfg.Profiles {
				output += fmt.Sprintf("- %s\n", profileName)
			}
		}

		return mcp.NewToolResultText(output), nil
	}
}

func listGradesHandler(cfg *config.Config, adapters map[string]spaggiari.Adapter) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		profileName, err := request.RequireString("profile")
		if err != nil {
			return mcp.NewToolResultError("profile is required"), nil
		}

		adapter, ok := adapters[profileName]
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("Unknown profile: %s", profileName)), nil
		}

		grades, err := adapter.Grades.List()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch grades: %v", err)), nil
		}

		limit := request.GetInt("limit", 0)
		if limit > 0 && limit < len(grades) {
			grades = grades[:limit]
		}

		data, err := json.MarshalIndent(grades, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal grades: %v", err)), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

func listAgendaHandler(cfg *config.Config, adapters map[string]spaggiari.Adapter) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		profileName, err := request.RequireString("profile")
		if err != nil {
			return mcp.NewToolResultError("profile is required"), nil
		}

		adapter, ok := adapters[profileName]
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("Unknown profile: %s", profileName)), nil
		}

		since := time.Now()
		sinceStr := request.GetString("since", "")
		if sinceStr != "" {
			parsed, err := time.Parse("2006-01-02", sinceStr)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid since date format: %v", err)), nil
			}
			since = parsed
		}

		until := time.Now().AddDate(0, 0, 30)
		untilStr := request.GetString("until", "")
		if untilStr != "" {
			parsed, err := time.Parse("2006-01-02", untilStr)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid until date format: %v", err)), nil
			}
			until = parsed
		}

		entries, err := adapter.Agenda.List(since, until)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch agenda: %v", err)), nil
		}

		limit := request.GetInt("limit", 0)
		if limit > 0 && limit < len(entries) {
			entries = entries[:limit]
		}

		data, err := json.MarshalIndent(entries, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal agenda: %v", err)), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

func listNoticeboardsHandler(cfg *config.Config, adapters map[string]spaggiari.Adapter) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		profileName, err := request.RequireString("profile")
		if err != nil {
			return mcp.NewToolResultError("profile is required"), nil
		}

		adapter, ok := adapters[profileName]
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("Unknown profile: %s", profileName)), nil
		}

		items, err := adapter.Noticeboards.List()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch noticeboards: %v", err)), nil
		}

		data, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal noticeboards: %v", err)), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

func getExampleConfig() string {
	return `profiles:
  older-kid:
    username: JOHNdoe123
    password: password123
  younger-kid:
    username: JANEDOE456
    password: password456
default_profile: older-kid
`
}
