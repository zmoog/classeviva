# MCP Server Development Guide

## Architecture

The MCP server is built using the [mcp-go](https://github.com/mark3labs/mcp-go) library and integrates with the existing Classeviva adapters.

### Components

1. **Configuration Layer** (`adapters/config/`)
   - `config.go`: Handles loading/saving student configurations
   - Stores credentials securely in `~/.classeviva/students.json`

2. **MCP Server** (`entrypoints/mcp/main.go`)
   - Implements the Model Context Protocol server
   - Provides 4 tools for LLM interaction
   - Reuses existing Spaggiari adapters for API calls

3. **Existing Adapters** (`adapters/spaggiari/`)
   - `identity.go`: Authentication and token management
   - `grades.go`: Grade retrieval
   - `agenda.go`: Agenda/homework retrieval
   - `noticeboards.go`: Noticeboard/announcement retrieval

### Flow

```
LLM (e.g., Claude)
    ↓
MCP Protocol (stdio)
    ↓
MCP Server (entrypoints/mcp/main.go)
    ↓
Student Config (adapters/config/config.go)
    ↓
Spaggiari Adapter (adapters/spaggiari/)
    ↓
Classeviva API (web.spaggiari.eu)
```

## Tool Implementations

### list_profiles
- Lists all configured profiles
- No API calls required
- Reads from configuration file

### list_grades
- Requires profile parameter
- Calls `adapter.Grades.List()`
- Optional limit parameter for pagination

### list_agenda
- Requires profile parameter
- Calls `adapter.Agenda.List(since, until)`
- Optional date range and limit parameters

### list_noticeboards
- Requires profile parameter
- Calls `adapter.Noticeboards.List()`
- Returns all noticeboards for the profile

## Testing

### Unit Tests
Run existing tests:
```bash
go test ./...
```

### Manual Testing
1. Create a test configuration:
   ```bash
   mkdir -p ~/.classeviva
   cp entrypoints/mcp/students.json.example ~/.classeviva/students.json
   # Edit with real credentials
   ```

2. Build and run:
   ```bash
   go build -o classeviva-mcp ./entrypoints/mcp/main.go
   ./classeviva-mcp
   ```

3. Test with MCP Inspector:
   ```bash
   npx @modelcontextprotocol/inspector classeviva-mcp
   ```

### Integration Testing with Claude
1. Configure Claude Desktop with the MCP server
2. Ask Claude to:
   - List students
   - Show grades for a specific student
   - Check upcoming homework
   - Review recent noticeboards

## Security Considerations

### Credential Storage
- Credentials stored in `~/.classeviva/config.yaml`
- Recommended file permissions: `0600` (owner read/write only)
- File created automatically with secure permissions

### API Authentication
- Uses existing identity management from Spaggiari adapter
- Tokens cached in `~/.classeviva/identity-*.json`
- Token refresh handled automatically

### No Credential Exposure
- Credentials never logged or returned in tool results
- Only profile names exposed to LLM
- Passwords never transmitted in MCP protocol

## Adding New Tools

To add a new MCP tool:

1. Identify the required Spaggiari adapter method
2. Create a handler function in `main.go`:
   ```go
   func myNewToolHandler(cfg config.Config, adapters map[string]spaggiari.Adapter) server.ToolHandlerFunc {
       return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
           // Implementation
       }
   }
   ```

3. Register the tool in the `run()` function:
   ```go
   s.AddTool(mcp.NewTool("my_new_tool",
       mcp.WithDescription("Description of the tool"),
       mcp.WithString("profile",
           mcp.Description("The profile name from list_profiles"),
           mcp.Required(),
       ),
   ), myNewToolHandler(cfg, profileAdapters))
   ```

## Troubleshooting

### Common Issues

1. **Profiles not found**
   - Ensure `~/.classeviva/config.yaml` exists and contains at least one student
   - Check file permissions

2. **Authentication failures**
   - Verify credentials are correct
   - Check if API is accessible from your location
   - Some cloud providers are blocked by Classeviva

3. **Server not starting**
   - Check for syntax errors: `go build ./entrypoints/mcp/main.go`
   - Verify dependencies: `go mod tidy`

4. **Claude not seeing the server**
   - Verify Claude Desktop configuration
   - Check server path is absolute
   - Restart Claude Desktop after config changes

### Debug Mode

Enable debug logging by setting the log level:
```go
log.SetLevel(log.TraceLevel)
```

## Future Enhancements

Possible additions:
1. Download noticeboard attachments
2. Filter grades by subject or date range
3. Calendar integration for agenda items
4. Attendance information
5. Teacher communications
6. Resource templates for structured queries

## Dependencies

- **mcp-go**: MCP protocol implementation
- **sirupsen/logrus**: Logging
- Existing Classeviva dependencies (see go.mod)
