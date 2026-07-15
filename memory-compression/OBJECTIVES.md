# Objectives — memory-compression

1. Compress long-horizon agent memory into dense, replayable state without losing provenance.
2. Reduce: Raw memory logs grow unbounded and slow every downstream read.
3. Maintain sub-second decision latency on the hot path where applicable.
4. Zero unhandled errors in production execution paths.
5. Full provenance for every action (hash-chained memory write).
