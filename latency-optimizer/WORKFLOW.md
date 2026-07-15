# Workflow — latency-optimizer

Receive Goal
    -> Understand Intent
    -> Build Plan (shared/planner)
    -> Acquire Context (shared/memory)
    -> Invoke MCP (shared/mcp, via tool-selector)
    -> Execute (Rust (hot path) + Go (control plane))
    -> Validate
    -> Store Memory (Postgres, hash-chained)
    -> Evaluate (shared/evaluation)
    -> Return Result
