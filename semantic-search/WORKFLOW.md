# Workflow — semantic-search

Receive Goal
    -> Understand Intent
    -> Build Plan (shared/planner)
    -> Acquire Context (shared/memory)
    -> Invoke MCP (shared/mcp, via tool-selector)
    -> Execute (Python (embeddings) + Rust (vector index) + Go (API))
    -> Validate
    -> Store Memory (Postgres, hash-chained)
    -> Evaluate (shared/evaluation)
    -> Return Result
