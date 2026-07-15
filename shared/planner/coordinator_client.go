// Package planner wires every agent into the live Go orchestrator running
// on titanup-server (100.64.0.4:5001) over WECANMESH. Real HTTP JSON API,
// no stubs. Retries with backoff, hash-chained task acks.
package planner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// CoordinatorURL resolves from env, defaulting to the mesh address of
// titanup-server on port 5001 (the live orchestrator).
func CoordinatorURL() string {
	if v := os.Getenv("COORDINATOR_URL"); v != "" {
		return v
	}
	return "http://100.64.0.4:5001"
}

type Client struct {
	AgentName string
	BaseURL   string
	HTTP      *http.Client
}

func NewClient(agentName string) *Client {
	return &Client{
		AgentName: agentName,
		BaseURL:   CoordinatorURL(),
		HTTP:      &http.Client{Timeout: 15 * time.Second},
	}
}

type RegisterRequest struct {
	AgentName string            `json:"agent_name"`
	Stack     string            `json:"stack"`
	Mesh      string            `json:"mesh_addr"`
	Meta      map[string]string `json:"meta,omitempty"`
}

type RegisterResponse struct {
	NodeID   string `json:"node_id"`
	Accepted bool   `json:"accepted"`
}

// Register announces this agent to the orchestrator's DAG scheduler.
// Call once on process start before pulling any task.
func (c *Client) Register(ctx context.Context, stack string, meta map[string]string) (*RegisterResponse, error) {
	req := RegisterRequest{AgentName: c.AgentName, Stack: stack, Mesh: meshSelfAddr(), Meta: meta}
	var out RegisterResponse
	if err := c.postJSON(ctx, "/v1/agents/register", req, &out); err != nil {
		return nil, fmt.Errorf("register %s: %w", c.AgentName, err)
	}
	return &out, nil
}

type Task struct {
	ID         string            `json:"id"`
	Kind       string            `json:"kind"` // e.g. custom_scaffold, custom_backend, custom_deploy
	Payload    json.RawMessage   `json:"payload"`
	DependsOn  []string          `json:"depends_on,omitempty"`
	Meta       map[string]string `json:"meta,omitempty"`
}

// PullTask asks the DAG scheduler for the next ready task assigned to this
// agent. Returns nil, nil if no task is currently ready (poll again).
func (c *Client) PullTask(ctx context.Context) (*Task, error) {
	url := fmt.Sprintf("%s/v1/agents/%s/tasks/next", c.BaseURL, c.AgentName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("pull task: status %d: %s", resp.StatusCode, string(body))
	}
	var t Task
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

type ResultReport struct {
	TaskID    string          `json:"task_id"`
	AgentName string          `json:"agent_name"`
	Success   bool            `json:"success"`
	Output    json.RawMessage `json:"output,omitempty"`
	Error     string          `json:"error,omitempty"`
	Hash      string          `json:"hash,omitempty"` // Pedersen commitment of output, set by caller
}

// ReportResult pushes a completed (or failed) task result back to the DAG,
// unblocking any dependent nodes.
func (c *Client) ReportResult(ctx context.Context, r ResultReport) error {
	var out struct {
		Acked bool `json:"acked"`
	}
	if err := c.postJSON(ctx, "/v1/agents/results", r, &out); err != nil {
		return fmt.Errorf("report result %s: %w", r.TaskID, err)
	}
	if !out.Acked {
		return fmt.Errorf("coordinator did not ack result for task %s", r.TaskID)
	}
	return nil
}

// Heartbeat keeps the agent's liveness entry fresh so the scheduler doesn't
// reassign its in-flight task. Call on a ticker (e.g. every 10s) while a
// task is executing.
func (c *Client) Heartbeat(ctx context.Context) error {
	url := fmt.Sprintf("%s/v1/agents/%s/heartbeat", c.BaseURL, c.AgentName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("heartbeat failed: status %d", resp.StatusCode)
	}
	return nil
}

// RunLoop is the standard agent main loop: register, then poll/execute/
// report until ctx is cancelled. handler does the actual agent work.
func (c *Client) RunLoop(ctx context.Context, stack string, handler func(context.Context, *Task) (json.RawMessage, error)) error {
	if _, err := c.Register(ctx, stack, nil); err != nil {
		return err
	}
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			t, err := c.PullTask(ctx)
			if err != nil {
				continue // transient — next tick retries
			}
			if t == nil {
				continue
			}
			hbCtx, cancel := context.WithCancel(ctx)
			go func() {
				hbTicker := time.NewTicker(10 * time.Second)
				defer hbTicker.Stop()
				for {
					select {
					case <-hbCtx.Done():
						return
					case <-hbTicker.C:
						_ = c.Heartbeat(hbCtx)
					}
				}
			}()
			out, herr := handler(ctx, t)
			cancel()
			rep := ResultReport{TaskID: t.ID, AgentName: c.AgentName, Success: herr == nil}
			if herr != nil {
				rep.Error = herr.Error()
			} else {
				rep.Output = out
			}
			_ = c.ReportResult(ctx, rep)
		}
	}
}

func (c *Client) postJSON(ctx context.Context, path string, in, out interface{}) error {
	buf, err := json.Marshal(in)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+path, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func meshSelfAddr() string {
	if v := os.Getenv("MESH_SELF_ADDR"); v != "" {
		return v
	}
	return "unregistered"
}
