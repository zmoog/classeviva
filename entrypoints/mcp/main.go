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

	configLoader := config.Loader{Path: homeDir}
	cfg, err := configLoader.Load()
	if err != nil {
		configPath, _ := configLoader.GetConfigPath()
		return fmt.Errorf("failed to load configuration: %w\n\nPlease create a configuration file at: %s\n\nExample format:\n%s",
			err, configPath, getExampleConfig())
	}

	studentAdapters := make(map[string]spaggiari.Adapter)
	for _, student := range cfg.Students {
		adapter, err := spaggiari.New(student.Username, student.Password, homeDir)
		if err != nil {
			return fmt.Errorf("failed to create adapter for student %s: %w", student.Name, err)
		}
		studentAdapters[student.ID] = adapter
	}

	s := server.NewMCPServer(
		"Classeviva MCP Server",
		"1.0.0",
		server.WithLogging(),
	)

	s.AddTool(mcp.NewTool("list_students",
		mcp.WithDescription("List all configured students"),
		mcp.WithString("format",
			mcp.Description("Output format: 'text' or 'json'"),
			mcp.DefaultString("json"),
		),
	), listStudentsHandler(cfg))

	s.AddTool(mcp.NewTool("list_grades",
		mcp.WithDescription("List grades for a specific student"),
		mcp.WithString("student_id",
			mcp.Description("The student ID from list_students"),
			mcp.Required(),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of grades to return (default: all)"),
		),
	), listGradesHandler(cfg, studentAdapters))

	s.AddTool(mcp.NewTool("list_agenda",
		mcp.WithDescription("List agenda items (homework, events) for a specific student"),
		mcp.WithString("student_id",
			mcp.Description("The student ID from list_students"),
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
	), listAgendaHandler(cfg, studentAdapters))

	s.AddTool(mcp.NewTool("list_noticeboards",
		mcp.WithDescription("List noticeboard items (announcements, circulars) for a specific student"),
		mcp.WithString("student_id",
			mcp.Description("The student ID from list_students"),
			mcp.Required(),
		),
	), listNoticeboardsHandler(cfg, studentAdapters))

	if err := server.ServeStdio(s); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

func listStudentsHandler(cfg config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		format := request.GetString("format", "json")

		students := make([]map[string]string, 0, len(cfg.Students))
		for _, student := range cfg.Students {
			students = append(students, map[string]string{
				"id":   student.ID,
				"name": student.Name,
			})
		}

		var output string
		if format == "json" {
			data, err := json.MarshalIndent(students, "", "  ")
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal students: %v", err)), nil
			}
			output = string(data)
		} else {
			output = "Available students:\n"
			for _, student := range cfg.Students {
				output += fmt.Sprintf("- ID: %s, Name: %s\n", student.ID, student.Name)
			}
		}

		return mcp.NewToolResultText(output), nil
	}
}

func listGradesHandler(cfg config.Config, adapters map[string]spaggiari.Adapter) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		studentID, err := request.RequireString("student_id")
		if err != nil {
			return mcp.NewToolResultError("student_id is required"), nil
		}

		adapter, ok := adapters[studentID]
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("Unknown student_id: %s", studentID)), nil
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

func listAgendaHandler(cfg config.Config, adapters map[string]spaggiari.Adapter) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		studentID, err := request.RequireString("student_id")
		if err != nil {
			return mcp.NewToolResultError("student_id is required"), nil
		}

		adapter, ok := adapters[studentID]
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("Unknown student_id: %s", studentID)), nil
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

func listNoticeboardsHandler(cfg config.Config, adapters map[string]spaggiari.Adapter) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		studentID, err := request.RequireString("student_id")
		if err != nil {
			return mcp.NewToolResultError("student_id is required"), nil
		}

		adapter, ok := adapters[studentID]
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("Unknown student_id: %s", studentID)), nil
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
	return `{
  "students": [
    {
      "id": "student1",
      "name": "John Doe",
      "username": "JOHNdoe123",
      "password": "password123"
    },
    {
      "id": "student2",
      "name": "Jane Doe",
      "username": "JANEDOE456",
      "password": "password456"
    }
  ]
}`
}
