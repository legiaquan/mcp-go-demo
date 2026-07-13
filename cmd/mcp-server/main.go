package main

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"mcp-go-demo/internal/config"
	"mcp-go-demo/internal/db"
	internalLogger "mcp-go-demo/internal/logger"
	"mcp-go-demo/internal/tools"
)

func main() {
	// Initialize logger
	internalLogger.Init()
	logger := internalLogger.Log
	logger.Info("Starting MCP Server")

	// Load configuration
	cfg := config.Load()

	// Initialize database
	database, err := db.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()
	logger.Info("Connected to database successfully")

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-go-demo",
		Version: "1.0.0",
	}, nil)

	// Register Task Tools
	taskTools := &tools.TaskTools{
		DB:     database,
		Logger: logger,
	}
	taskTools.Register(server)

	// Register Weather Tools
	weatherTools := &tools.WeatherTools{
		Logger: logger,
	}
	weatherTools.Register(server)

	// Start standard I/O transport
	logger.Info("Starting stdio transport")
	transport := mcp.NewStdioTransport()
	if err := server.Connect(context.Background(), transport); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	logger.Info("Server stopped")
}
