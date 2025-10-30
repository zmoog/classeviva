# Classeviva MCP Server

An MCP (Model Context Protocol) server implementation for Classeviva, enabling LLMs like Claude to interact with the Classeviva school portal API.

## Overview

The MCP server provides tools to:
- List configured student profiles
- Retrieve student grades
- Access student agenda (homework, events)
- View noticeboard items (announcements, circulars)

## Configuration

The MCP server requires a configuration file at `~/.classeviva/config.yaml` with student profile credentials.

Create the configuration file:

```bash
mkdir -p ~/.classeviva
cat > ~/.classeviva/config.yaml << 'EOF'
profiles:
  older-kid:
    username: JOHNdoe123
    password: password123
  younger-kid:
    username: JANEDOE456
    password: password456
default_profile: older-kid
EOF
chmod 600 ~/.classeviva/config.yaml
```

Or manually create `~/.classeviva/config.yaml` with your credentials:

```yaml
profiles:
  older-kid:
    username: JOHNdoe123
    password: password123
  younger-kid:
    username: JANEDOE456
    password: password456
default_profile: older-kid
```

**Security Note**: The configuration file should have permissions `0600` (owner read/write only) to protect sensitive credentials.

## Installation

### Quick Install

Run the installation script:

```bash
cd /path/to/classeviva
./entrypoints/mcp/install.sh
```

This will:
1. Build the MCP server binary
2. Install it to `/usr/local/bin` or `~/.local/bin`
3. Create a configuration template at `~/.classeviva/config.yaml`

### Manual Build

Alternatively, build the MCP server manually:

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

### list_profiles

List all configured student profiles.

**Parameters:**
- `format` (optional): Output format - `"text"` or `"json"` (default: `"json"`)

**Example:**
```
Use list_profiles to see available student profiles
```

### list_grades

List grades for a specific student profile.

**Parameters:**
- `profile` (required): Profile name from list_profiles
- `limit` (optional): Maximum number of grades to return

**Example:**
```
Show me the grades for older-kid
```

### list_agenda

List agenda items (homework, events) for a specific student profile.

**Parameters:**
- `profile` (required): Profile name from list_profiles
- `since` (optional): Start date in YYYY-MM-DD format (default: today)
- `until` (optional): End date in YYYY-MM-DD format (default: 30 days from now)
- `limit` (optional): Maximum number of items to return

**Example:**
```
Show me the agenda for younger-kid for the next week
```

### list_noticeboards

List noticeboard items (announcements, circulars) for a specific student profile.

**Parameters:**
- `profile` (required): Profile name from list_profiles

**Example:**
```
Show me the latest noticeboards for older-kid
```

## Example Workflows

With the MCP server, you can ask Claude to:

1. **Analyze grades**: "Compare the grades between my two student profiles and identify areas where each student is excelling or needs support"

2. **Check homework**: "What homework assignments are due this week for older-kid?"

3. **Review announcements**: "Summarize the recent noticeboard announcements for both student profiles"

4. **Track progress**: "Show me all the recent grades for younger-kid in Mathematics and Science"

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
