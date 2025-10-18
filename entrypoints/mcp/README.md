# Classeviva MCP Server

An MCP (Model Context Protocol) server implementation for Classeviva, enabling LLMs like Claude to interact with the Classeviva school portal API.

## Overview

The MCP server provides tools to:
- List configured students
- Retrieve student grades
- Access student agenda (homework, events)
- View noticeboard items (announcements, circulars)

## Configuration

The MCP server requires a configuration file at `~/.classeviva/students.json` with student credentials:

```json
{
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
}
```

**Security Note**: The configuration file is created with permissions `0600` (owner read/write only) to protect sensitive credentials.

## Building

Build the MCP server:

```bash
go build -o classeviva-mcp ./entrypoints/mcp/main.go
```

## Using with Claude Desktop

To use the MCP server with Claude Desktop, add it to your Claude configuration file:

### macOS
Edit `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "classeviva": {
      "command": "/path/to/classeviva-mcp"
    }
  }
}
```

### Windows
Edit `%APPDATA%\Claude\claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "classeviva": {
      "command": "C:\\path\\to\\classeviva-mcp.exe"
    }
  }
}
```

### Linux
Edit `~/.config/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "classeviva": {
      "command": "/path/to/classeviva-mcp"
    }
  }
}
```

After updating the configuration, restart Claude Desktop.

## Available Tools

### list_students

List all configured students.

**Parameters:**
- `format` (optional): Output format - `"text"` or `"json"` (default: `"json"`)

**Example:**
```
Use list_students to see available students
```

### list_grades

List grades for a specific student.

**Parameters:**
- `student_id` (required): Student ID from list_students
- `limit` (optional): Maximum number of grades to return

**Example:**
```
Show me the grades for student1
```

### list_agenda

List agenda items (homework, events) for a specific student.

**Parameters:**
- `student_id` (required): Student ID from list_students
- `since` (optional): Start date in YYYY-MM-DD format (default: today)
- `until` (optional): End date in YYYY-MM-DD format (default: 30 days from now)
- `limit` (optional): Maximum number of items to return

**Example:**
```
Show me the agenda for student2 for the next week
```

### list_noticeboards

List noticeboard items (announcements, circulars) for a specific student.

**Parameters:**
- `student_id` (required): Student ID from list_students

**Example:**
```
Show me the latest noticeboards for student1
```

## Example Workflows

With the MCP server, you can ask Claude to:

1. **Analyze grades**: "Compare the grades between my two students and identify areas where each student is excelling or needs support"

2. **Check homework**: "What homework assignments are due this week for student1?"

3. **Review announcements**: "Summarize the recent noticeboard announcements for both students"

4. **Track progress**: "Show me all the recent grades for student2 in Mathematics and Science"

## Troubleshooting

### Configuration file not found

If you see an error about the configuration file not existing, create it at the specified path with your student credentials.

### Authentication errors

If you encounter authentication errors:
1. Verify your credentials are correct
2. Ensure you're not calling the API from a blocked location (some cloud providers are blocked by Classeviva)
3. Check that identity tokens are not expired (they are cached in `~/.classeviva/identity.json`)

### Connection issues

If the server fails to connect:
1. Check your internet connection
2. Verify that `https://web.spaggiari.eu` is accessible
3. Review the server logs for specific error messages
