# Architecture

25-agent fleet, each independently deployable, coordinated by a Go DAG scheduler (coordinator agent) with Rust/Tokio workers.
Shared substrate: Postgres event fabric, Pedersen/Ristretto255 ZK provenance chain, WECANMESH sovereign WireGuard mesh, local GGUF inference (zero external API dependency by default).
Polyglot stack per workload: Go (orchestration/API), Rust (performance-critical/concurrency), Zig (WASM scoring kernels, deterministic compute), Python (ML/statistics/NLP), C++ (numeric cores, GPU-adjacent muscle).
