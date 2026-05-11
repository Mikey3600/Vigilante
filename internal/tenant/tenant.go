package tenant

import (
	"context"
)

type contextKey string

const tenantKey contextKey = "tenant_id"

// NewContext returns a new context injected with a tenant identifier.
func NewContext(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantKey, tenantID)
}

// FromContext extracts the bounded tenant from a given Context execution path.
func FromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(tenantKey).(string)
	return val, ok
}
