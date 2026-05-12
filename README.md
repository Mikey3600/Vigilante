# Vigilante

![Go Version](https://img.shields.io/badge/go-1.22-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green)
![Status](https://img.shields.io/badge/status-active-brightgreen)

Self-hostable backend observability platform. Vigilante ingests logs and metrics from your microservices, detects anomalies automatically, and uses an LLM to generate incident summaries with root cause analysis and fix suggestions — all from your own infrastructure.

No third-party SaaS. No data leaves your environment unless you call the Groq API.

---

## What it does

- Ingests logs and metrics over REST and gRPC from any microservice
- Stores time-series data in TimescaleDB for efficient range queries
- Runs an anomaly engine that detects spikes, drops, and error patterns
- Sends anomalies to Groq (LLaMA 3.1) for human-readable incident summaries with root cause and fix suggestions
- Auto-saves AI analysis results to the dashboard in real time
- Serves a real-time dashboard at `localhost:3000`
- Supports multiple tenants with JWT-based authentication

---

## Requirements

- Go 1.22+
- Docker Desktop (for PostgreSQL + TimescaleDB)
- A free Groq API key from [console.groq.com](https://console.groq.com)

---

## Quickstart

1. **Clone the repository**

```bash
   git clone https://github.com/Mikey3600/Vigilante.git
   cd Vigilante
```

2. **Configure environment**

```bash
   cp .env.example .env
```

   Open `.env` and fill in the required values:

   | Variable | Description |
   |---|---|
   | `DATABASE_URL` | PostgreSQL connection string (e.g. `postgres://vigilante:vigilante@localhost:5432/vigilante`) |
   | `JWT_SECRET` | A long random string used to sign tokens |
   | `GROQ_API_KEY` | Your API key from console.groq.com |

3. **Start the database**

```bash
   docker-compose up -d postgres
```

4. **Build the binary**

```bash
   go build -o vigilante.exe ./cmd/vigilante
```

5. **Run database migrations**

```bash
   ./vigilante.exe migrate
```

6. **Create the first tenant and get a token**

```bash
   ./vigilante.exe setup
```

   This prints a JWT token to stdout. Copy it — you'll need it to log into the dashboard. Token expires in 30 days. Re-run `setup` to get a new one.

7. **Start the server**

```bash
   ./vigilante.exe serve
```

8. **Open the dashboard**

   Navigate to [http://localhost:3000](http://localhost:3000) and paste the JWT token from step 6 when prompted.

---

## Sending data to Vigilante

Once the server is running, send metrics and logs from any service using the JWT token from setup:

**Send a metric:**
```powershell
Invoke-WebRequest -Uri "http://localhost:3000/api/v1/metrics" `
  -Method POST `
  -Headers @{ "Authorization" = "Bearer YOUR_TOKEN"; "Content-Type" = "application/json" } `
  -Body '{"service_id":"11111111-1111-1111-1111-111111111111","metrics":[{"metric_name":"latency_ms","value":350}]}'
```

**Send a log:**
```powershell
Invoke-WebRequest -Uri "http://localhost:3000/api/v1/logs" `
  -Method POST `
  -Headers @{ "Authorization" = "Bearer YOUR_TOKEN"; "Content-Type" = "application/json" } `
  -Body '{"service_id":"11111111-1111-1111-1111-111111111111","logs":[{"level":"error","message":"database connection timeout after 30s","time":"2026-05-12T10:00:00Z"}]}'
```

**Trigger AI analysis (auto-saves to dashboard):**
```powershell
Invoke-WebRequest -Uri "http://localhost:3000/api/v1/analyze" `
  -Method POST `
  -Headers @{ "Authorization" = "Bearer YOUR_TOKEN"; "Content-Type" = "application/json" } `
  -Body '{"service_id":"11111111-1111-1111-1111-111111111111","anomaly_type":"latency_spike","description":"Service latency spiked to 8000ms"}'
```

---

## CLI Commands

| Command | Description |
|---|---|
| `serve` | Start the REST API, gRPC server, and dashboard |
| `migrate` | Apply all pending database migrations |
| `setup` | Provision the initial tenant and print a JWT token |
| `status` | Show server health, DB connectivity, and active tenant count |
| `ingest` | Manually push a log or metric payload from a file or stdin |

---

## Architecture

```mermaid
graph LR
    A[Microservice 1] --> C[REST API :3000]
    B[Microservice 2] --> D[gRPC Server :50051]
    C --> E[(TimescaleDB)]
    D --> E
    E --> F[Anomaly Engine]
    F --> G[AI Analyzer - Groq]
    G --> H[Dashboard]
    E --> H
```

---

## Tech Stack

| Component | Technology |
|---|---|
| Language | Go 1.22 |
| HTTP framework | Gin |
| Database | PostgreSQL + TimescaleDB |
| Authentication | JWT |
| AI inference | Groq API (LLaMA 3.1) |
| Service ingestion | gRPC |
| CLI | Cobra |
| Local environment | Docker |

---

## License

[MIT](LICENSE) — © Vigilante Contributors