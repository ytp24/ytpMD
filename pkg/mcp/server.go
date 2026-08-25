package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Server coordinates the Model Context Protocol stdio JSON-RPC 2.0 transport.
type Server struct {
	reader io.Reader
	writer io.Writer
}

// NewServer initializes a new MCP Server instance.
func NewServer(reader io.Reader, writer io.Writer) *Server {
	return &Server{
		reader: reader,
		writer: writer,
	}
}

// StartServer starts the MCP stdio loop on standard input and standard output.
func StartServer() {
	// Diagnostic logs strictly to stderr
	fmt.Fprintf(os.Stderr, "[ytpMD MCP] Server starting in stdio mode (Version 3.2.0)...\n")
	server := NewServer(os.Stdin, os.Stdout)
	if err := server.Listen(context.Background()); err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "[ytpMD MCP] Server stopped with error: %v\n", err)
	}
}

// Listen processes incoming JSON-RPC 2.0 messages from the configured reader.
func (s *Server) Listen(ctx context.Context) error {
	scanner := bufio.NewScanner(s.reader)
	// Buffer size up to 10MB for large JSON-RPC messages
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req JSONRPCRequest
		if err := json.Unmarshal(line, &req); err != nil {
			s.sendError(nil, CodeParseError, "Parse error: invalid JSON", nil)
			continue
		}

		s.handleRequest(ctx, req)
	}

	return scanner.Err()
}

func (s *Server) handleRequest(ctx context.Context, req JSONRPCRequest) {
	switch req.Method {
	case "initialize":
		result := InitializeResult{
			ProtocolVersion: "2024-11-05",
			ServerInfo: ServerInfo{
				Name:    "ytpmd-mcp",
				Version: "3.2.0",
			},
			Capabilities: ServerCapabilities{
				Tools: &ToolsCapability{
					ListChanged: false,
				},
				Resources: &ResourcesCapability{
					Subscribe:   false,
					ListChanged: false,
				},
				Prompts: &PromptsCapability{
					ListChanged: false,
				},
			},
			Instructions: "ytpMD MCP Server provides tools to convert PDF documents and technical books into chapter-split, agent-ready Markdown notes with YAML frontmatter, breadcrumb navigation, and AGENTS.md manifests.",
		}
		s.sendResponse(req.ID, result)

	case "notifications/initialized", "initialized":
		// Client notification after initialization, no response needed for notifications without ID
		if req.ID != nil {
			s.sendResponse(req.ID, map[string]interface{}{})
		}

	case "ping":
		s.sendResponse(req.ID, map[string]interface{}{})

	case "tools/list":
		tools := GetRegisteredTools()
		s.sendResponse(req.ID, ListToolsResult{Tools: tools})

	case "tools/call":
		var params CallToolParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			s.sendError(req.ID, CodeInvalidParams, "Invalid params for tools/call", err.Error())
			return
		}

		result := ExecuteTool(ctx, params.Name, params.Arguments)
		s.sendResponse(req.ID, result)

	case "resources/list":
		s.sendResponse(req.ID, map[string]interface{}{"resources": []interface{}{}})

	case "prompts/list":
		s.sendResponse(req.ID, map[string]interface{}{"prompts": []interface{}{}})

	default:
		// If request has an ID, send MethodNotFound error
		if req.ID != nil {
			s.sendError(req.ID, CodeMethodNotFound, fmt.Sprintf("Method not found: '%s'", req.Method), nil)
		}
	}
}

func (s *Server) sendResponse(id interface{}, result interface{}) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	s.writeJSON(resp)
}

func (s *Server) sendError(id interface{}, code int, message string, data interface{}) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	s.writeJSON(resp)
}

func (s *Server) writeJSON(v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ytpMD MCP] Failed to marshal response: %v\n", err)
		return
	}
	_, _ = s.writer.Write(append(data, '\n'))
}
