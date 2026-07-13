# MCP Go Demo Server

This is a demonstration of building a [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) Server in Go. It provides tools for both interacting with an external API (National Weather Service) and a PostgreSQL database (Task Management).

## Architecture

This project follows the **Standard Go Project Layout** combined with principles of Clean Architecture to ensure maintainability and testability:

```text
mcp-go-demo/
├── cmd/
│   └── mcp-server/
│       └── main.go            # Entry point: Initialize DB, setup MCP Server and register Tools
├── internal/
│   ├── config/                # Environment configuration
│   ├── logger/                # Structured logging (slog) setup
│   ├── db/                    # PostgreSQL connection and migrations
│   ├── models/                # Core Data Models (Tasks, NWS Responses, MCP Inputs)
│   ├── service/               # External Services / 3rd Party APIs (Weather API logic)
│   └── tools/                 # MCP Handlers / Controllers (Task & Weather Tools)
├── docker-compose.yml         # Defines MCP Server + PostgreSQL DB + MCP Inspector UI
├── Dockerfile                 # Multi-stage Docker build for the Go Server
└── README.md
```

## Tools Provided

The MCP server exposes the following tools to the LLM:

### 1. Task Management (PostgreSQL Database)
- `add_task`: Add a new task (fields: title, status)
- `update_task`: Update a task title or status
- `delete_task`: Delete a task by ID
- `get_task`: Retrieve a single task by ID
- `list_tasks`: List all tasks in the database

### 2. Weather (External API)
- `get_forecast`: Retrieve a detailed weather forecast for a latitude/longitude pair
- `get_alerts`: Retrieve active weather alerts for a US state

## Running Locally via Docker Compose

We use `docker compose` to orchestrate the PostgreSQL database, the Go MCP server, and the MCP Inspector frontend for easy testing.

1. Bring up the stack:
   ```bash
   docker compose up -d --build
   ```

2. Open the **MCP Inspector** in your browser:
   Check the terminal logs or run `docker logs -f mcp-go-demo-mcp-inspector-1` to find the proxy connection URL (e.g., `http://localhost:6274/?MCP_PROXY_AUTH_TOKEN=...`).

3. Testing your tools:
   In the MCP Inspector UI, you can query available tools, invoke them, and view responses.

## Logging & Tracing

This project uses `log/slog` to write JSON-structured logs to `stderr`. `stdout` is strictly reserved for the MCP JSON-RPC protocol. Every tool invocation logs its lifecycle (`Started` / `Completed`) with a unique `trace_id` and tracks execution latency to facilitate deep debugging.
