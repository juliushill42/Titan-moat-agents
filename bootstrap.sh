#!/usr/bin/env bash
# Titan agent fleet bootstrap — real inference + coordinator wiring, no stubs.
# Run as: bash bootstrap.sh   (never ./bootstrap.sh — exFAT has no exec bit)
set -euo pipefail

echo "[titan-agents] bootstrap start"

command -v go >/dev/null || { echo "FATAL: go not found"; exit 1; }
command -v cargo >/dev/null || { echo "FATAL: cargo not found"; exit 1; }
command -v python3 >/dev/null || { echo "FATAL: python3 not found"; exit 1; }
command -v curl >/dev/null || { echo "FATAL: curl not found"; exit 1; }

if [ ! -f .env ]; then
  cp .env.example .env
  echo "[titan-agents] .env created from .env.example"
fi
set -a; source .env; set +a

LLAMA_SERVER_URL="${LLAMA_SERVER_URL:-http://100.64.0.7:8091}"
COORDINATOR_URL="${COORDINATOR_URL:-http://100.64.0.4:5001}"

echo "[titan-agents] verifying live inference at ${LLAMA_SERVER_URL}"
if curl -fsS -m 5 -X POST "${LLAMA_SERVER_URL}/completion" \
    -H "Content-Type: application/json" \
    -d '{"prompt":"ping","temperature":0,"seed":42,"top_k":1,"n_predict":4}' >/tmp/llama_check.json 2>/tmp/llama_check.err; then
  echo "[titan-agents] LLAMA_SERVER_URL is live: ${LLAMA_SERVER_URL}"
else
  echo "[titan-agents] WARNING: LLAMA_SERVER_URL unreachable (${LLAMA_SERVER_URL}) — see /tmp/llama_check.err"
  echo "[titan-agents]          agents will fail closed on inference until justus-sovereign:8091 is up"
fi

echo "[titan-agents] verifying coordinator at ${COORDINATOR_URL}"
if curl -fsS -m 5 "${COORDINATOR_URL}/v1/health" >/tmp/coord_check.json 2>/tmp/coord_check.err; then
  echo "[titan-agents] coordinator is live: ${COORDINATOR_URL}"
else
  echo "[titan-agents] WARNING: coordinator unreachable (${COORDINATOR_URL}) — see /tmp/coord_check.err"
  echo "[titan-agents]          agents will retry registration via backoff in RunLoop"
fi

echo "[titan-agents] initializing Go modules + wiring inference/coordinator per agent"
for d in agents/*/ ; do
  agent="$(basename "$d")"
  if [ ! -f "${d}go.mod" ]; then
    (cd "$d" && go mod init "titanuniversalai.com/agents/${agent}" >/dev/null 2>&1 || true)
  fi

  mkdir -p "${d}connectors"

  cat > "${d}connectors/inference.go" <<EOF
// Auto-wired by bootstrap.sh — do not hand-edit; regenerate via bootstrap.sh.
package connectors

import (
	"context"

	"titanuniversalai.com/shared/llm"
)

var llmRouter = llm.NewRouterFromEnv()

// Infer runs a sovereign-first completion for the ${agent} agent: local
// GGUF (Qwen2.5 via justus-sovereign:8091) first, Gemini 2.5 Flash only if
// ALLOW_GEMINI_FALLBACK=true and local inference fails.
func Infer(ctx context.Context, prompt string) (string, error) {
	res, err := llmRouter.Complete(ctx, llm.CompletionRequest{Prompt: prompt})
	if err != nil {
		return "", err
	}
	return res.Text, nil
}
EOF

  cat > "${d}connectors/coordinator.go" <<EOF
// Auto-wired by bootstrap.sh — do not hand-edit; regenerate via bootstrap.sh.
package connectors

import (
	"context"
	"encoding/json"

	"titanuniversalai.com/shared/planner"
)

var Coordinator = planner.NewClient("${agent}")

// Start registers this agent with the live orchestrator on titanup-server
// (:5001) and runs the pull/execute/report loop until ctx is cancelled.
func Start(ctx context.Context, stack string, handler func(context.Context, *planner.Task) (json.RawMessage, error)) error {
	return Coordinator.RunLoop(ctx, stack, handler)
}
EOF

done

echo "[titan-agents] bootstrap complete — ${LLAMA_SERVER_URL} / ${COORDINATOR_URL} wired into $(ls agents | wc -l) agents"
