# Objectives — semantic-search

1. Serve low-latency semantic search over all Titan repos, docs, and memory stores.
2. Reduce: Keyword search misses intent; existing search tools require external API calls (non-sovereign).
3. Maintain sub-second decision latency on the hot path where applicable.
4. Zero unhandled errors in production execution paths.
5. Full provenance for every action (hash-chained memory write).
