// Package llm routes every agent's inference call to the sovereign local
// GGUF server first (LLAMA_SERVER_URL, Qwen2.5 via llama.cpp), and only
// falls back to Gemini 2.5 Flash when explicitly enabled — never silently.
// Go owns the API/routing layer per the polyglot stack convention; Python
// and C++ own the actual model math (llama.cpp, PyTorch/TensorFlow).
package llm

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

type Backend string

const (
	BackendLocalGGUF Backend = "local_gguf"
	BackendGemini    Backend = "gemini_2_5_flash"
)

type CompletionRequest struct {
	Prompt      string
	Temperature float64 // sovereign default: 0
	Seed        int     // sovereign default: 42
	TopK        int     // sovereign default: 1
	MaxTokens   int
}

type CompletionResult struct {
	Text    string
	Backend Backend
}

type Router struct {
	LlamaURL       string // e.g. http://100.64.0.7:8091 (justus-sovereign llama-server)
	GeminiAPIKey   string
	GeminiModel    string // "gemini-2.5-flash" — override via GEMINI_MODEL if Google renames/versions it
	AllowFallback  bool   // must be explicitly true; sovereign-first by default
	HTTP           *http.Client
}

func NewRouterFromEnv() *Router {
	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = "gemini-2.5-flash"
	}
	return &Router{
		LlamaURL:      envOr("LLAMA_SERVER_URL", "http://100.64.0.7:8091"),
		GeminiAPIKey:  os.Getenv("GEMINI_API_KEY"),
		GeminiModel:   model,
		AllowFallback: os.Getenv("ALLOW_GEMINI_FALLBACK") == "true",
		HTTP:          &http.Client{Timeout: 60 * time.Second},
	}
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// Complete tries local GGUF first. If it errors and AllowFallback is set
// and a GeminiAPIKey is present, it falls back to Gemini 2.5 Flash and
// tags the result so callers/telemetry know sovereignty was broken for
// that call.
func (r *Router) Complete(ctx context.Context, req CompletionRequest) (*CompletionResult, error) {
	if req.Temperature == 0 && req.Seed == 0 && req.TopK == 0 {
		req.Temperature, req.Seed, req.TopK = 0, 42, 1 // WECANMESH deterministic default
	}
	text, err := r.completeLocal(ctx, req)
	if err == nil {
		return &CompletionResult{Text: text, Backend: BackendLocalGGUF}, nil
	}
	if !r.AllowFallback || r.GeminiAPIKey == "" {
		return nil, fmt.Errorf("local GGUF failed and Gemini fallback disabled/unconfigured: %w", err)
	}
	text, gerr := r.completeGemini(ctx, req)
	if gerr != nil {
		return nil, fmt.Errorf("local GGUF failed (%v) and Gemini fallback failed (%v)", err, gerr)
	}
	return &CompletionResult{Text: text, Backend: BackendGemini}, nil
}

// --- local llama.cpp server (/completion) ---

type llamaCompletionReq struct {
	Prompt      string  `json:"prompt"`
	Temperature float64 `json:"temperature"`
	Seed        int     `json:"seed"`
	TopK        int     `json:"top_k"`
	NPredict    int     `json:"n_predict"`
}

type llamaCompletionResp struct {
	Content string `json:"content"`
}

func (r *Router) completeLocal(ctx context.Context, req CompletionRequest) (string, error) {
	body, _ := json.Marshal(llamaCompletionReq{
		Prompt:      req.Prompt,
		Temperature: req.Temperature,
		Seed:        req.Seed,
		TopK:        req.TopK,
		NPredict:    orDefault(req.MaxTokens, 512),
	})
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, r.LlamaURL+"/completion", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := r.HTTP.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("llama-server unreachable at %s: %w", r.LlamaURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("llama-server status %d: %s", resp.StatusCode, string(b))
	}
	var out llamaCompletionResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.Content, nil
}

// --- Gemini 2.5 Flash (generativelanguage.googleapis.com) ---

type geminiPart struct {
	Text string `json:"text"`
}
type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}
type geminiReq struct {
	Contents []geminiContent `json:"contents"`
}
type geminiCandidate struct {
	Content geminiContent `json:"content"`
}
type geminiResp struct {
	Candidates []geminiCandidate `json:"candidates"`
}

func (r *Router) completeGemini(ctx context.Context, req CompletionRequest) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		r.GeminiModel, r.GeminiAPIKey)
	body, _ := json.Marshal(geminiReq{Contents: []geminiContent{{Parts: []geminiPart{{Text: req.Prompt}}}}})
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := r.HTTP.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gemini status %d: %s", resp.StatusCode, string(b))
	}
	var out geminiResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if len(out.Candidates) == 0 || len(out.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini returned no candidates")
	}
	return out.Candidates[0].Content.Parts[0].Text, nil
}

func orDefault(v, def int) int {
	if v == 0 {
		return def
	}
	return v
}
