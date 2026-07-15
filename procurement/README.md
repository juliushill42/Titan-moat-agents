# Procurement

## Mission
Track vendor/tool spend and contract terms, and flag renewal or leverage windows.

## Real-world Bottleneck
Recurring costs and contract terms are scattered and easy to miss.

## Target Users
Titan Universal AI internal agent fleet; JU (Julius Cameron Hill) as sole operator/approver.

## Inputs
Event streams, task graph goals, MCP tool results, prior agent memory.

## Outputs
Structured results written to shared/memory (Postgres), telemetry events, and downstream task graph updates.

## Decision Process
Scan billing/contract sources -> Go normalizes into ledger -> flags renewal windows and cost anomalies.

## Planning Strategy
Goal decomposition against shared/planner; no agent plans in isolation — every plan node registers with the Go coordinator DAG.

## Memory Strategy
Short-term: in-process context. Long-term: Postgres via shared/memory, compressed by the memory-compression agent, hash-chained (Pedersen/Ristretto255) for provenance.

## MCP Integrations
Registered through shared/mcp; tool selection delegated to tool-selector agent.

## External Tools
None by default — sovereign-first. External calls only via explicit MCP grant.

## Evaluation Metrics
Task success rate, latency (p50/p99), token/compute cost, human override rate.

## Security Model
BSL 1.1 licensing, patent poison-pill compliance enforced by governance agent. No training data leaves Titan infra.

## Failure Recovery
Self-heal kill/respawn pattern (TAAG-style); failed tasks requeue with backoff; incident-commander escalates on repeat failure.

## Observability
Structured telemetry to shared/telemetry; traced spans for latency-optimizer consumption.
