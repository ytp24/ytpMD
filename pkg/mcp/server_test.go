package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestMCPServer_Initialize(t *testing.T) {
	in := bytes.NewBufferString(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}` + "\n")
	var out bytes.Buffer

	server := NewServer(in, &out)
	err := server.Listen(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var resp JSONRPCResponse
	if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v, raw: %s", err, out.String())
	}

	if resp.ID.(float64) != 1 {
		t.Errorf("expected id 1, got %v", resp.ID)
	}

	resultMap, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map result, got %T", resp.Result)
	}

	if resultMap["protocolVersion"] != "2024-11-05" {
		t.Errorf("expected protocolVersion 2024-11-05, got %v", resultMap["protocolVersion"])
	}
}

func TestMCPServer_ListTools(t *testing.T) {
	in := bytes.NewBufferString(`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}` + "\n")
	var out bytes.Buffer

	server := NewServer(in, &out)
	_ = server.Listen(context.Background())

	var resp JSONRPCResponse
	if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	resultMap, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map result, got %T", resp.Result)
	}

	tools, ok := resultMap["tools"].([]interface{})
	if !ok || len(tools) < 4 {
		t.Errorf("expected at least 4 registered tools, got %d", len(tools))
	}
}

func TestMCPServer_InspectPDF_InvalidPath(t *testing.T) {
	req := `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"inspect_pdf","arguments":{"path":"non_existent_file.pdf"}}}` + "\n"
	in := bytes.NewBufferString(req)
	var out bytes.Buffer

	server := NewServer(in, &out)
	_ = server.Listen(context.Background())

	var resp JSONRPCResponse
	if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	resultMap, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map result, got %T", resp.Result)
	}

	if resultMap["isError"] != true {
		t.Errorf("expected isError true for non existent file, got %v", resultMap["isError"])
	}
}

func TestMCPServer_GetManifest_Missing(t *testing.T) {
	req := `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"get_manifest","arguments":{"folder_path":"/tmp/non_existent_dir_12345"}}}` + "\n"
	in := bytes.NewBufferString(req)
	var out bytes.Buffer

	server := NewServer(in, &out)
	_ = server.Listen(context.Background())

	var resp JSONRPCResponse
	_ = json.Unmarshal(out.Bytes(), &resp)

	resultMap := resp.Result.(map[string]interface{})
	if resultMap["isError"] != true {
		t.Errorf("expected isError true for non-existent manifest directory")
	}
}

func TestMCPServer_MethodNotFound(t *testing.T) {
	req := `{"jsonrpc":"2.0","id":5,"method":"unknown/method","params":{}}` + "\n"
	in := bytes.NewBufferString(req)
	var out bytes.Buffer

	server := NewServer(in, &out)
	_ = server.Listen(context.Background())

	var resp JSONRPCResponse
	_ = json.Unmarshal(out.Bytes(), &resp)

	if resp.Error == nil {
		t.Fatalf("expected error object, got nil")
	}

	if resp.Error.Code != CodeMethodNotFound {
		t.Errorf("expected error code %d, got %d", CodeMethodNotFound, resp.Error.Code)
	}
	if !strings.Contains(resp.Error.Message, "Method not found") {
		t.Errorf("expected method not found message, got %s", resp.Error.Message)
	}
}
