package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/ItsJooL/valkey-mcp-server/internal/client"
	"github.com/ItsJooL/valkey-mcp-server/internal/registry"
	"github.com/ItsJooL/valkey-mcp-server/internal/tools"
	"github.com/ItsJooL/valkey-mcp-server/internal/types"
)

func main() {
	log.SetOutput(os.Stderr)

	urlFlag := flag.String("url", "", "Valkey connection URL (e.g., valkey://localhost:6379)")
	passwordFlag := flag.String("password", "", "Valkey password")
	dbFlag := flag.Int("db", 0, "Valkey database number (0-15)")
	flag.Parse()

	valkeyURL := *urlFlag
	if valkeyURL == "" {
		valkeyURL = os.Getenv("VALKEY_URL")
	}
	if valkeyURL == "" {
		valkeyURL = "valkey://localhost:6379"
	}

	password := *passwordFlag
	if password == "" {
		password = os.Getenv("VALKEY_PASSWORD")
	}

	db := *dbFlag
	if dbEnv := os.Getenv("VALKEY_DB"); dbEnv != "" {
		fmt.Sscanf(dbEnv, "%d", &db)
	}

	url, err := types.NewValkeyURL(valkeyURL)
	if err != nil {
		log.Fatalf("Invalid Valkey URL: %v", err)
	}

	dbIndex, err := types.NewDBIndex(db)
	if err != nil {
		log.Fatalf("Invalid database index: %v", err)
	}

	ctx := context.Background()

	valkeyClient, err := client.New(ctx, client.Config{
		URL:      url,
		Password: password,
		DB:       dbIndex,
	})
	if err != nil {
		log.Fatalf("Failed to create Valkey client: %v", err)
	}
	defer valkeyClient.Close()
	toolRegistry := registry.NewToolRegistry()
	tools.RegisterAll(toolRegistry, valkeyClient)

	log.Printf("Valkey MCP Server started")
	log.Printf("Connected to: %s (DB: %d)", url, dbIndex.Int())
	log.Printf("Available tools: %d", toolRegistry.Count())

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "valkey-mcp-server",
		Version: "1.0.0",
	}, nil)

	if err := toolRegistry.RegisterWithMCP(server); err != nil {
		log.Fatalf("Failed to register tools with MCP: %v", err)
	}

	log.Println("MCP server initialized, starting stdio transport...")

	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}
