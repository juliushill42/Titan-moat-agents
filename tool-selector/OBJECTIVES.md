# Objectives — tool-selector

1. Route an agent's intent to the correct tool/MCP server with minimal latency and no wasted calls.
2. Reduce: Agents call the wrong tool or too many tools, burning tokens and time.
3. Maintain sub-second decision latency on the hot path where applicable.
4. Zero unhandled errors in production execution paths.
5. Full provenance for every action (hash-chained memory write).
