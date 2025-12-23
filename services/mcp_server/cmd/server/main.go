package main

import (
	"log"
	"mcp_server/internal/config"
	"mcp_server/internal/tools"
	"net/http"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	cfg := config.Load()

	// ser := service.NewChangeService(cfg)

	// changeHandler := tools.NewChangeHandler(ser)
	// scheduleHandler := tools.NewSchedulesHandler(ser)

	tgHandler := tools.NewTGHandler(cfg)

	mcpServer := server.NewMCPServer(
		"Доступ к расписанию и изменениям расписания",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	scheduleTools := tools.NewScheduleTools()
	// mcpServer.AddTool(scheduleTools.GetChanges(), changeHandler.HandleChangeRequest)
	// mcpServer.AddTool(scheduleTools.GetCurrentDate(), scheduleHandler.CurrentDateHandler)
	// mcpServer.AddTool(scheduleTools.GetCurrentStudyPeriod(), scheduleHandler.CurrentStudyPeriodHandler)
	// mcpServer.AddTool(scheduleTools.GetGroupSchedule(), scheduleHandler.GetGroupSchedule)
	mcpServer.AddTool(scheduleTools.SendChanges(), tgHandler.SendChanges)

	http.Handle("/mcp", server.NewStreamableHTTPServer(mcpServer))
	log.Printf("MCP chedule server starting on :%s", cfg.ServerPort)
	http.ListenAndServe(":"+cfg.ServerPort, nil)
}
