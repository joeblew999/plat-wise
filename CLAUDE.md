# plat-wise - Wise Banking Platform

Go client for the [Wise API](https://docs.wise.com/api-reference) with CLI, MCP server, and web GUI.

## xplat

This project uses [xplat](https://github.com/joeblew999/xplat) for task running and process orchestration. See the [xplat CLAUDE.md](../xplat/CLAUDE.md) for conventions on Taskfile structure, variable naming, and process-compose patterns.

## Important

**Keep Taskfile.yml in sync with the API** - When adding new API endpoints or commands, always update the Taskfile with corresponding tasks.

## Structure

```
plat-wise/
├── client.go         # HTTP client with services
├── oauth.go          # OAuth 2.0 authentication
├── errors.go         # API error types
├── types.go          # Common types (Currency, Money, Timestamp)
├── profiles.go       # Profiles API
├── quotes.go         # Quotes API
├── recipients.go     # Recipients API
├── transfers.go      # Transfers API
├── rates.go          # Exchange rates API
├── balances.go       # Balances API
├── commands/         # Shared business logic (DRY)
│   └── commands.go
├── cmd/
│   ├── wise-cli/     # CLI tool
│   ├── wise-mcp/     # MCP server for Claude
│   └── wise-server/  # Web GUI with Via
├── Taskfile.yml      # Task commands
├── API.md            # API coverage documentation
└── README.md         # Quick start
```

## Architecture

- **API client** (`*.go`) - Reusable library for Wise API
- **commands** - Shared business logic, returns data structures
- **wise-cli** - CLI that formats command output for terminal
- **wise-mcp** - MCP server that formats output for Claude
- **wise-server** - Web GUI with Via framework (SSE for live updates)

All tools use the `commands` package to avoid duplication (DRY).

## Authentication

### API Token (Simple)
```bash
export WISE_API_TOKEN=your-token-here
```

### OAuth 2.0 (Multi-user / Partners)
```bash
export WISE_CLIENT_ID=your-client-id
export WISE_CLIENT_SECRET=your-client-secret
export WISE_REDIRECT_URL=http://localhost:8080/oauth/callback  # optional
```

OAuth flow:
1. User redirects to `wise.com/oauth/authorize`
2. User grants access
3. Wise redirects back with authorization code
4. Exchange code for access token
5. Token auto-refreshes (12 hour expiry)

## Wise API Endpoints

### Profiles
- `GET /v2/profiles` - List all profiles
- `GET /v2/profiles/{id}` - Get profile by ID

### Balances
- `GET /v4/profiles/{id}/balances` - List balances (requires `types=STANDARD`)
- `GET /v1/profiles/{id}/balance-statements/{balanceId}/statement.json` - Get statements

### Exchange Rates
- `GET /v1/rates` - Get rates (public, no auth needed)
- `GET /v1/rates?from=...&to=...&group=day` - Get rate history

### Quotes
- `POST /v2/quotes` - Create quote
- `GET /v2/quotes/{id}` - Get quote

### Recipients
- `POST /v1/accounts` - Create recipient
- `GET /v1/accounts` - List recipients

### Transfers
- `POST /v1/transfers` - Create transfer
- `GET /v1/transfers/{id}` - Get transfer

## Tasks

```bash
# CLI commands
task rates         # Get exchange rates
task profiles      # List profiles
task balances      # Show balances
task statements    # Transaction history
task quote         # Get currency quote
task rate-history  # Get historical rates

# MCP server
task mcp           # Run MCP server
task mcp-build     # Build MCP binary
task mcp-config    # Print Claude Desktop config

# Web GUI
task serve         # Start web dashboard (port 8080)
task serve-build   # Build web server binary

# Build/Test
task build         # Build all binaries
task test          # Run tests
task clean         # Remove built binaries

# Debug
task debug         # Print env vars
```

## CLI Help

```bash
./wise-cli -h                        # General help
./wise-cli -cmd help rate-history    # Help for specific command
```

## MCP Server Tools

- `wise_rates` - Get exchange rates between currency pairs
- `wise_profiles` - List all Wise profiles
- `wise_balances` - Show account balances
- `wise_statements` - Get transaction history
- `wise_quote` - Get currency conversion quotes
- `wise_rate_history` - Get historical exchange rates

## Web GUI Features

- Profiles list
- Live account balances
- Exchange rate table
- Currency conversion quotes
- Transaction statements
- Rate history with charts
- OAuth login flow (when configured)
- Real-time updates via SSE (Via framework)

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `WISE_API_TOKEN` | Yes* | Personal API token |
| `WISE_CLIENT_ID` | Yes* | OAuth client ID |
| `WISE_CLIENT_SECRET` | Yes* | OAuth client secret |
| `WISE_REDIRECT_URL` | No | OAuth redirect (default: localhost) |
| `WISE_SANDBOX` | No | Set to "true" for sandbox |

*Either API token OR OAuth credentials required.

## Wise API Notes

- Access tokens expire after 12 hours
- Refresh tokens should be stored securely
- Some endpoints require OAuth (not personal tokens) in EU/UK due to PSD2
- Rate limits apply - check response headers
- Sandbox: `api.sandbox.transferwise.tech`
- Production: `api.wise.com`

## Links

- [Wise API Reference](https://docs.wise.com/api-reference)
- [Auth & Security Guide](https://docs.wise.com/guides/developer/auth-and-security)
- [OAuth User Tokens](https://docs.wise.com/api-reference/user-tokens)
