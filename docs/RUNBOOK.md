# Operations Runbook

## How to Deploy
1. Update `docker-compose.yml` to specify persistent volumes if in Prod.
2. Apply Kubernetes primitives via `kubectl apply -f deploy/k8s/`
3. Make sure the Postgres database has TimescaleDB natively installed.

## Managing Tenants
* Tenants currently require explicit database insertion.
* Run: `INSERT INTO tenants (name) VALUES ('Acme Corp');`
* Proceed to create users with the provided `tenant_id` hash bounds.

## Debugging Anomaly Engine Loops
1. Make sure background `checkMetrics()` has a valid target DB context.
2. If `checkMetrics()` halts, examine `log_entries` sizes. It may require chunked query processing.

## Key Rotation
1. Update `JWT_SECRET` in `.env` / ConfigMap.
2. Restart the deployment using `kubectl rollout restart deployment vigilante`.
3. Inform clients that an invalidation sequence kicked off.

## Emergency Steps
* To halt data ingestion safely: Scale the pod horizontally to `0`. Postgres guarantees durability for whatever processed entirely prior to termination.
