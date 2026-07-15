// Package mcp implements a minimal JSON-RPC 2.0 client over the Model
// Context Protocol so any Go agent can list and call tools on a registered
// MCP server without duplicating transport logic.
package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type Server struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

type Client struct {
	server Server
	http   *http.Client
	id     int64
}

func NewClient(s Server) *Client {
	return &Client{server: s, http: &http.Client{Timeout: 30 * time.Second}}
}

type rpcRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int64       `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

func (c *Client) call(ctx context.Context, method string, params interface{}, out interface{}) error {
	id := atomic.AddInt64(&c.id, 1)
	body, err := json.Marshal(rpcRequest{JSONRPC: "2.0", ID: id, Method: method, Params: params})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.server.URL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("mcp %s call %s: %w", c.server.Name, method, err)
	}
	defer resp.Body.Close()
	var rr rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		return fmt.Errorf("mcp %s decode: %w", c.server.Name, err)
	}
	if rr.Error != nil {
		return fmt.Errorf("mcp %s error %d: %s", c.server.Name, rr.Error.Code, rr.Error.Message)
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(rr.Result, out)
}

// ListTools returns the tools this MCP server exposes.
func (c *Client) ListTools(ctx context.Context) ([]Tool, error) {
	var out struct {
		Tools []Tool `json:"tools"`
	}
	if err := c.call(ctx, "tools/list", nil, &out); err != nil {
		return nil, err
	}
	return out.Tools, nil
}

// CallTool invokes a named tool with arguments and returns the raw result
// content. The tool-selector agent decides which server/tool to hit; this
// client just executes the dispatch.
func (c *Client) CallTool(ctx context.Context, name string, args map[string]interface{}) (json.RawMessage, error) {
	var out struct {
		Content json.RawMessage `json:"content"`
		IsError bool            `json:"isError"`
	}
	params := map[string]interface{}{"name": name, "arguments": args}
	if err := c.call(ctx, "tools/call", params, &out); err != nil {
		return nil, err
	}
	if out.IsError {
		return nil, fmt.Errorf("mcp tool %s reported error: %s", name, string(out.Content))
	}
	return out.Content, nil
}

// Registry holds every MCP server an agent fleet is allowed to reach,
// scoped least-privilege per agent by the governance agent.
type Registry struct {
	servers map[string]Server
}

func NewRegistry() *Registry { return &Registry{servers: map[string]Server{}} }

func (r *Registry) Add(s Server) { r.servers[s.Name] = s }

func (r *Registry) Client(name string) (*Client, error) {
	s, ok := r.servers[name]
	if !ok {
		return nil, fmt.Errorf("mcp server %q not registered", name)
	}
	return NewClient(s), nil
}
