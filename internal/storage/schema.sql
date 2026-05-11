CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS log_entries (
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    service_id UUID REFERENCES services(id),
    level VARCHAR(50),
    message TEXT,
    metadata JSONB
);
SELECT create_hypertable('log_entries', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS metric_points (
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    service_id UUID REFERENCES services(id),
    metric_name VARCHAR(255) NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    labels JSONB
);
SELECT create_hypertable('metric_points', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS anomalies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id UUID REFERENCES services(id),
    detected_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    anomaly_type VARCHAR(255) NOT NULL,
    description TEXT,
    root_cause_summary TEXT,
    likely_cause TEXT,
    suggested_fix TEXT
);

CREATE TABLE IF NOT EXISTS alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id),
    anomaly_id UUID REFERENCES anomalies(id),
    channel VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
