package server

import (
	"context"
	"mcp_server/internal/config"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	httpSrv *http.Server
	McpSrv *server.MCPServer
	cfg *config.Config
}

func NewMCPServer(cfg *config.Config) *MCPServer {
	mcpServer := server.NewMCPServer(
		"Доступ к расписанию и изменениям расписания",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	httpSrv := &http.Server{
		Addr: ":" + cfg.ServerPort,
	}


	return &MCPServer{
		httpSrv: httpSrv,
		McpSrv: mcpServer,
		cfg: cfg,
	}
}

func (s *MCPServer) RegisterTool(tool mcp.Tool, handler mcp.StructuredToolHandlerFunc[any, any]) {
	s.McpSrv.AddTool(tool, mcp.NewStructuredToolHandler(handler))
}

func (s *MCPServer) Start() error {
	r := http.NewServeMux()
	r.Handle("/mcp", server.NewStreamableHTTPServer(s.McpSrv))

	s.httpSrv.Handler = r

	return s.httpSrv.ListenAndServe()
}

func (s *MCPServer) Stop(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}
