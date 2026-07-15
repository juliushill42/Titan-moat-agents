# Evaluation — process-mining

- Correctness: task output validated against schema + spot-check sampling
- Efficiency: latency p50/p99 tracked via shared/telemetry
- Cost: token/compute spend tracked via cost-optimizer
- Latency: hot-path profiled by latency-optimizer
- Reliability: success rate over rolling 7-day window
- Safety: zero unauthorized irreversible actions
- Tool Effectiveness: per-tool success rate feeding tool-selector ranking
- Continuous Improvement Plan: prompt-optimizer runs scheduled A/B cycles on this agent's prompt
