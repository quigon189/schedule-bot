package main

import (
	"log"
	"mcp_server/internal/config"
	"mcp_server/internal/service"
	"mcp_server/internal/tools"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	cfg := config.Load()

	handler := tools.NewChangeHandler(
		service.NewChangeService(cfg),
	)

	mcpServer := server.NewMCPServer(
		"Доступ к расписанию и изменениям расписания",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	getChangesTool := mcp.NewTool(
		"get_schedule_changes",
		mcp.WithDescription("Позволяет получить ссылку на изображение(я), в которой указаны изменения в расписании для всех груп на заданную в формате ISO 8601 дату"),
		mcp.WithString("date",
			mcp.Description("Дата в формате ISO 8601 YYYY-MM-DD. Если не задана, то возвращает изменения на текущую дату"),
		),
	)

	mcpServer.AddTool(getChangesTool, handler.HandleChangeRequest)

	http.Handle("/mcp", server.NewStreamableHTTPServer(mcpServer))
	log.Printf("MCP chedule server starting on :%s", cfg.ServerPort)
	http.ListenAndServe(":"+cfg.ServerPort, nil)
}
