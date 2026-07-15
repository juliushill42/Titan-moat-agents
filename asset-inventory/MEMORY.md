# Memory — asset-inventory

## Short-Term Context
In-process working set for the active task graph node.

## Long-Term Knowledge
Persisted to Postgres via shared/memory; compressed periodically by memory-compression agent.

## Retrieved Documents
Pulled via semantic-search agent when historical context is needed.

## Decisions
Every decision logged with timestamp, inputs, output, and Pedersen commitment hash.

## Lessons Learned
Fed to the mentor agent for cross-repo pattern transfer.

## Feedback Loop
Evaluation scores loop back into prompt-optimizer for this agent's prompt variants.
