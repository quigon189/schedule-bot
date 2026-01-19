package main

import (
	"context"
	"log"
	"mcp_server/internal/config"
	"mcp_server/internal/server"
	"mcp_server/internal/tools"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	cfg := config.Load()

	srv := server.NewMCPServer(cfg)

	tgTools := tools.NewScheduleTools()
	tgHandlers := tools.NewTGHandler(cfg)

	srv.McpSrv.AddTool(tgTools.SendChanges(), mcp.NewStructuredToolHandler(tgHandlers.SendChanges))

	go func() {
		log.Printf("Server started on port %s", cfg.ServerPort)
		log.Println(srv.Start())
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		log.Fatalf("Failed to stop server: %v", err)
		return
	}

	log.Printf("Server stoped")
}
