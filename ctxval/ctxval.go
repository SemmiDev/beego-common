// Package ctxval provides type-safe helpers for storing and retrieving
// values from a context.Context using Go generics.
//
// Usage:
//
//	type ctxKey string
//	const userKey ctxKey = "user"
//
//	ctx = ctxval.Set(ctx, userKey, currentUser)
//	user, ok := ctxval.Get[User](ctx, userKey)
package ctxval

import "context"

// Set stores value in ctx under key and returns the updated context.
//
//	ctx = ctxval.Set(ctx, userKey, currentUser)
func Set[T any](ctx context.Context, key, value any) context.Context {
	return context.WithValue(ctx, key, value)
}

// Get retrieves a value of type T from ctx by key.
// Returns the zero value of T and false if the key is absent or the value is
// not of type T.
//
//	user, ok := ctxval.Get[User](ctx, userKey)
func Get[T any](ctx context.Context, key any) (T, bool) {
	v := ctx.Value(key)
	if v == nil {
		var zero T
		return zero, false
	}
	typed, ok := v.(T)
	return typed, ok
}

// MustGet retrieves a value of type T from ctx by key and panics when the
// key is absent or the stored value is not of type T.
// Use only in code paths where the key is guaranteed to exist.
//
//	user := ctxval.MustGet[User](ctx, userKey)
func MustGet[T any](ctx context.Context, key any) T {
	v, ok := Get[T](ctx, key)
	if !ok {
		panic("ctxval: key not found or wrong type in context")
	}
	return v
}

// GetOrDefault retrieves a value of type T from ctx by key.
// Returns defaultVal when the key is absent or the stored value is not of
// type T.
//
//	lang := ctxval.GetOrDefault[string](ctx, langKey, "id")
func GetOrDefault[T any](ctx context.Context, key any, defaultVal T) T {
	v, ok := Get[T](ctx, key)
	if !ok {
		return defaultVal
	}
	return v
}
