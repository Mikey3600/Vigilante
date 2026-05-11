# Architecture Deep-Dive

## Data Flow Diagram

1. **Ingestion**: Log payloads or metrics hit `internal/api` endpoints or the gRPC bound ports in `internal/grpc`. 
2. **Buffer & Persist**: Data validates context and is asynchronously inserted using `pgx` into TimescaleDB hyper tables.
3. **Detection**: A goroutine periodically scans tables for anomalies (sudden spikes, latency out of standard deviation bounds).
4. **Resolution**: If triggered, `internal/ai` polls the recent log sequences and queries `gemini-1.5-flash` to get an inferred `RootCauseReport`.
5. **Routing**: Analyzed reports are persisted into the `anomalies` table and dispatch channels (WS, Slack, Email) are triggered through `internal/alert`.

## Tech Choices Explained

- **Go 1.22**: Chosen for extreme high-concurrency capabilities mapping thousands of ingestion points with low memory footprints.
- **TimescaleDB**: Used because log retention and metrics querying fundamentally requires high insert efficiency and sliding window selects. Traditional SQL would fail at scale.
- **Gemini API**: We utilize google/generative-ai-go and 1.5-flash for its extreme speed and wide context window limits, perfect for bulk log aggregation.
- **Gin Web Framework**: Clean HTTP routing and JWT middlewares.

## Anomaly Algorithm
Currently built around statistical variance. By scanning a sliding window (e.g. past 15 minutes) vs historical trend over a span of days, we map acceptable latency outputs. A violation triggers the alert mechanism.

## Multi-Tenancy Strategy
Handled at the HTTP request boundary via `auth.JWTMiddleware()`. When an auth token is decoded, `tenant_id` is assigned directly to the `gin.Context`.
A `tenant_id` bounds all SQL selects enforcing isolation natively.

## Scalability Notes
Vigilante targets horizontal elasticity. 
- HTTP nodes can be placed behind a standard load balancer.
- Background anomaly engine requires a distributed lock or partition strategy across multiple nodes if operating in massive compute realms (e.g., locking via Redis or Postgres advisory locks).
