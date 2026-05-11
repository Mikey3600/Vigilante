# API Reference

All requests require an `Authorization` header containing the multi-tenant JWT context:
`Authorization: Bearer <token>`

---

## 1. POST `/api/v1/logs`

Ingests a bulk payload of logs for a particular service.

**Request Body**
```json
{
  "service_id": "ee0f2271-9f20-4a87-8dcb-62b14421d019",
  "logs": [
    {
      "time": "2026-05-11T10:40:00Z",
      "level": "ERROR",
      "message": "Connection refused to upstream Redis",
      "metadata": {"source": "auth-service"}
    }
  ]
}
```

**Response (200 OK)**
```json
{
  "status": "accepted"
}
```

---

## 2. POST `/api/v1/metrics`

Ingests time-series structured data.

**Request Body**
```json
{
  "service_id": "ee0f2271-9f20-4a87-8dcb-62b14421d019",
  "metrics": [
    {
      "time": "2026-05-11T10:40:00Z",
      "metric_name": "http.latency",
      "value": 150.45,
      "labels": {"path": "/checkout"}
    }
  ]
}
```

**Response (200 OK)**
```json
{
  "status": "accepted"
}
```

---

## 3. GET `/health`

Health validation for probes. (Unauthenticated)

**Response (200 OK)**
```json
{
  "status": "ok"
}
```
