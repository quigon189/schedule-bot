package main

import (
	"mcp_server/internal/config"
	"mcp_server/internal/server"
	"mcp_server/internal/tools"

	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	cfg := config.Load()

	srv := server.NewMCPServer(cfg)

	tgTools := tools.NewScheduleTools()
	tgHandlers := tools.NewTGHandler(cfg)

	srv.McpSrv.AddTool(tgTools.SendChanges(), mcp.NewStructuredToolHandler(tgHandlers.SendChanges))
}
