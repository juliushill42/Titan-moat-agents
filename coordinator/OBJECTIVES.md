# Objectives — coordinator

1. Central task graph scheduler that assigns work across all 25 agents with dependency awareness.
2. Reduce: Agents step on each other or stall waiting on undeclared dependencies.
3. Maintain sub-second decision latency on the hot path where applicable.
4. Zero unhandled errors in production execution paths.
5. Full provenance for every action (hash-chained memory write).
