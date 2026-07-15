# agents.titanuniversalai.com

25-agent sovereign fleet for Titan Universal AI LLC. Owner: Julius Cameron Hill (JU).

## Stack
Go, Rust, Zig, Python, C++, TypeScript/Tailwind (dashboard), Tokio, Postgres, Lua (where UGC-adjacent).
Sovereign-first: local GGUF inference, zero mandatory external API dependency, Pedersen/Ristretto255 ZK provenance chain.

## Agents (25)
- **workflow-discovery** — Go (orchestrator) + Python (pattern mining, sklearn/HMM)
- **process-mining** — Python (PM4Py-style engine, custom) + Go (API)
- **data-quality** — Go (pipeline) + Python (validation rules/statistics)
- **api-builder** — Go
- **semantic-search** — Python (embeddings) + Rust (vector index) + Go (API)
- **prompt-optimizer** — Python
- **tool-selector** — Go
- **cost-optimizer** — Go (control loop) + Python (forecasting)
- **latency-optimizer** — Rust (hot path) + Go (control plane)
- **memory-compression** — Rust + Zig (WASM scoring kernel)
- **coordinator** — Go + Rust/Tokio
- **simulation** — Rust + C++ (physics/numerics core)
- **scenario-planning** — Python
- **mentor** — Python
- **documentation** — Go (extraction) + Python (generation)
- **incident-commander** — Go
- **meeting-intelligence** — Python
- **procurement** — Go
- **grant-writer** — Python
- **asset-inventory** — Go
- **risk-intelligence** — Python + Go
- **governance** — Go
- **digital-twin** — Rust + C++
- **experiment-runner** — Python + Go
- **innovation-scout** — Python

## Shared
- **planner** — Go
- **router** — Go
- **memory** — Rust + Postgres
- **prompts** — Python
- **evaluation** — Python
- **telemetry** — Go
- **policies** — Go
- **schemas** — Protobuf + Go codegen
- **mcp** — Go
- **utilities** — Go

## Bootstrap
```
bash bootstrap.sh
```

## License
Titan Proprietary License 1.0 — see LICENSE.
