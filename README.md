# Vigilante

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)
![Status](https://img.shields.io/badge/status-active-brightgreen)

A self-hostable Backend Observability & Incident Intelligence Platform written in Go. Vigilante ingests logs and metrics from your services, detects anomalies automatically, and uses AI to generate human-readable root-cause analysis — like having an on-call SRE that never sleeps.

> Built as a portfolio project demonstrating Go backend engineering, distributed systems, TimescaleDB, gRPC, JWT auth, and AI integration.

---

## What it does

- Ingests logs and metrics from any service via REST or gRPC
- Stores time-series data in TimescaleDB (PostgreSQL extension)
- Detects anomalies using statistical analysis (rolling average + stddev)
- Calls AI (Groq LLaMA 3.1) to generate incident summaries, root causes, and fixes
- Displays everything on a real-time dashboard
- Goes offline gracefully — dashboard freezes when backend dies, reconnects when it comes back

---

## Architecture

```mermaid
graph LR
    subgraph Your Services
        A[Microservice 1]
        B[Microservice 2]
    end

    subgraph Vigilante Platform
        C[REST API :3000]
        D[gRPC Server :50051]
        E[(TimescaleDB)]
        F[Anomaly Engine]
        G[AI Analyzer - Groq]
        H[Dashboard]
    end

    A --> C
    B --> D
    C --> E
    D --> E
    F  E
    F --> G
    G --> H
    E --> H
```

---

## Quickstart

### Prerequisites
- Go 1.22+
- Docker Desktop
- A free Groq API key from https://console.groq.com

### 1. Clone the repo
```bash
git clone https://github.com/Mikey3600/Vigilante.git
cd Vigilante
```

### 2. Set up environment
```bash
cp .env.example .env
```

Open `.env` and fill in:
```
DATABASE_URL=postgres://vigilante:vigilante@localhost:5432/vigilante?sslmode=disable
JWT_SECRET=your-secret-key-here
GROQ_API_KEY=your-groq-api-key-here
GRPC_PORT=50051
HTTP_PORT=3000
ENV=development
```

### 3. Start the database
```bash
docker-compose up -d postgres
```

### 4. Build the binary
```bash
go build -o vigilante.exe ./cmd/vigilante   # Windows
go build -o vigilante ./cmd/vigilante       # Mac/Linux
```

### 5. Run migrations
```bash
./vigilante.exe migrate
```

### 6. Run setup (creates tenant, service, and prints your JWT token)
```bash
./vigilante.exe setup
```

Copy the JWT token it prints — you'll need it for the dashboard.

### 7. Start the server
```bash
./vigilante.exe serve
```

### 8. Open the dashboard
Go to http://localhost:3000 and paste your JWT token.

---

## Sending telemetry

Once the server is running, send metrics from any service:

```powershell
$headers = @{ Authorization = "Bearer YOUR_TOKEN" }

# Send a metric
Invoke-WebRequest -Uri "http://localhost:3000/api/v1/metrics" `
  -Method POST -Headers $headers -ContentType "application/json" `
  -Body '{"service_id":"11111111-1111-1111-1111-111111111111","metrics":[{"metric_name":"latency","value":350}]}'

# Send a log
Invoke-WebRequest -Uri "http://localhost:3000/api/v1/logs" `
  -Method POST -Headers $headers -ContentType "application/json" `
  -Body '{"service_id":"11111111-1111-1111-1111-111111111111","logs":[{"level":"ERROR","message":"DB connection timeout"}]}'
```

---

## CLI Commands

| Command | Description |
|---|---|
| `vigilante serve` | Start HTTP and gRPC servers |
| `vigilante migrate` | Run database migrations |
| `vigilante setup` | Create default tenant/service and print JWT token |
| `vigilante ingest --file=logs.json` | Send a log file to the API |
| `vigilante status` | Check server health |

---

## Environment Variables

| Variable | Description | Required |
|---|---|---|
| `DATABASE_URL` | PostgreSQL connection string | ✅ |
| `JWT_SECRET` | Secret key for JWT signing | ✅ |
| `GROQ_API_KEY` | Groq API key for AI analysis | ✅ |
| `HTTP_PORT` | HTTP server port (default 3000) | ❌ |
| `GRPC_PORT` | gRPC server port (default 50051) | ❌ |
| `ENV` | Environment (development/production) | ❌ |

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.22 |
| Web Framework | Gin |
| Database | PostgreSQL + TimescaleDB |
| Auth | JWT (golang-jwt) |
| AI | Groq API (LLaMA 3.1 8B) |
| Transport | REST + gRPC |
| CLI | Cobra |
| Frontend | HTML + Chart.js + WebSocket |
| Containers | Docker + docker-compose |

---

## Project Structure
```
vigilante/
├── cmd/vigilante/     # CLI entry point (serve, migrate, setup, status)
├── internal/
│   ├── ai/            # Groq AI integration
│   ├── anomaly/       # Anomaly detection engine
│   ├── api/           # Gin HTTP handlers
│   ├── auth/          # JWT middleware
│   ├── grpc/          # gRPC server
│   ├── ingestion/     # Log and metric processing
│   ├── storage/       # PostgreSQL queries
│   └── tenant/        # Multi-tenant context
├── dashboard/         # Frontend HTML/JS
├── deploy/            # Docker and Kubernetes manifests
└── docs/              # Architecture and API docs
```

---

## License

MIT