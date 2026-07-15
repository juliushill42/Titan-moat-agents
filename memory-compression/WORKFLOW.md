# Workflow — memory-compression

Receive Goal
    -> Understand Intent
    -> Build Plan (shared/planner)
    -> Acquire Context (shared/memory)
    -> Invoke MCP (shared/mcp, via tool-selector)
    -> Execute (Rust + Zig (WASM scoring kernel))
    -> Validate
    -> Store Memory (Postgres, hash-chained)
    -> Evaluate (shared/evaluation)
    -> Return Result
