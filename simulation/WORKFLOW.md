# Workflow — simulation

Receive Goal
    -> Understand Intent
    -> Build Plan (shared/planner)
    -> Acquire Context (shared/memory)
    -> Invoke MCP (shared/mcp, via tool-selector)
    -> Execute (Rust + C++ (physics/numerics core))
    -> Validate
    -> Store Memory (Postgres, hash-chained)
    -> Evaluate (shared/evaluation)
    -> Return Result
