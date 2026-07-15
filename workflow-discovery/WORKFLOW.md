# Workflow — workflow-discovery

Receive Goal
    -> Understand Intent
    -> Build Plan (shared/planner)
    -> Acquire Context (shared/memory)
    -> Invoke MCP (shared/mcp, via tool-selector)
    -> Execute (Go (orchestrator) + Python (pattern mining, sklearn/HMM))
    -> Validate
    -> Store Memory (Postgres, hash-chained)
    -> Evaluate (shared/evaluation)
    -> Return Result
