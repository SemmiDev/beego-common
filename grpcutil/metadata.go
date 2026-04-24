package grpcutil

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	MetaKeyRequestID    = "x-request-id"
	MetaKeyAuthorization = "authorization"
	MetaKeyUserID       = "x-user-id"
	MetaKeyTraceID      = "x-trace-id"
	MetaKeyTenantID     = "x-tenant-id"
)

// MetaGet returns the first value for a metadata key from the incoming context.
// Keys are lowercased automatically.
func MetaGet(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	vals := md.Get(strings.ToLower(key))
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

// MetaGetAll returns all values for a metadata key from the incoming context.
func MetaGetAll(ctx context.Context, key string) []string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	return md.Get(strings.ToLower(key))
}

// MetaSet attaches key-value pairs to outgoing metadata on the context.
func MetaSet(ctx context.Context, kv ...string) context.Context {
	md := metadata.Pairs(kv...)
	return metadata.NewOutgoingContext(ctx, md)
}

// MetaMerge merges the provided key-value pairs into any existing outgoing
// metadata on the context.
func MetaMerge(ctx context.Context, kv ...string) context.Context {
	existing, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return MetaSet(ctx, kv...)
	}
	merged := existing.Copy()
	extra := metadata.Pairs(kv...)
	for k, v := range extra {
		merged[k] = append(merged[k], v...)
	}
	return metadata.NewOutgoingContext(ctx, merged)
}

// MetaGetIncoming extracts full incoming metadata map from context.
func MetaGetIncoming(ctx context.Context) metadata.MD {
	md, _ := metadata.FromIncomingContext(ctx)
	return md
}

// ExtractBearerToken returns the bearer token from the authorization metadata,
// stripping the "Bearer " prefix (case-insensitive).
func ExtractBearerToken(ctx context.Context) string {
	raw := MetaGet(ctx, MetaKeyAuthorization)
	if raw == "" {
		return ""
	}
	lower := strings.ToLower(raw)
	if strings.HasPrefix(lower, "bearer ") {
		return raw[7:]
	}
	return raw
}

// ExtractRequestID returns the request ID from incoming metadata.
func ExtractRequestID(ctx context.Context) string {
	return MetaGet(ctx, MetaKeyRequestID)
}

// ExtractUserID returns the user ID from incoming metadata.
func ExtractUserID(ctx context.Context) string {
	return MetaGet(ctx, MetaKeyUserID)
}

// ExtractTenantID returns the tenant ID from incoming metadata.
func ExtractTenantID(ctx context.Context) string {
	return MetaGet(ctx, MetaKeyTenantID)
}

// InjectOutgoing appends standard propagation headers to outgoing calls.
// Useful when forwarding metadata from an incoming RPC to a downstream RPC.
func InjectOutgoing(ctx context.Context) context.Context {
	pairs := []string{}
	for _, key := range []string{
		MetaKeyRequestID,
		MetaKeyTraceID,
		MetaKeyTenantID,
		MetaKeyUserID,
	} {
		if val := MetaGet(ctx, key); val != "" {
			pairs = append(pairs, key, val)
		}
	}
	if len(pairs) == 0 {
		return ctx
	}
	return MetaMerge(ctx, pairs...)
}
