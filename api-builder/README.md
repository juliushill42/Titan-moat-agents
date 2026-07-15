# Api Builder

## Mission
Generate, version, and harden internal/external REST and gRPC APIs from a spec, with zero hand-written boilerplate.

## Real-world Bottleneck
Every new service reinvents auth, validation, and error handling.

## Target Users
Titan Universal AI internal agent fleet; JU (Julius Cameron Hill) as sole operator/approver.

## Inputs
Event streams, task graph goals, MCP tool results, prior agent memory.

## Outputs
Structured results written to shared/memory (Postgres), telemetry events, and downstream task graph updates.

## Decision Process
Read OpenAPI/proto spec -> codegen handlers, middleware, validation -> attach Pedersen audit hooks -> deploy.

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
