# MCP Integration — workflow-discovery

| Field | Value |
|---|---|
| Partner | Registered per-deployment via shared/mcp |
| Purpose | Tool access for workflow discovery execution |
| Authentication | Scoped credential, least-privilege |
| Available Tools | Resolved dynamically via tool-selector |
| Planning Strategy | Cost/latency/success-rate ranked dispatch |
| Fallback Strategy | Retry with backoff -> requeue -> escalate to incident-commander |
| Rate Limits | Enforced by cost-optimizer control loop |
| Error Handling | Fail loud, log, never silently degrade output quality |
| Telemetry | All calls emit spans to shared/telemetry |
| Example Workflow | See WORKFLOW.md |
