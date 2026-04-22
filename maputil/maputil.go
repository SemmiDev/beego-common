// Package maputil provides generic utility functions for working with maps.
// All functions use Go 1.18+ generics.
//
// Usage:
//
//	keys := maputil.Keys(m)
//	filtered := maputil.Filter(m, func(k string, v int) bool { return v > 0 })
//	merged := maputil.Merge(defaults, overrides)
package maputil

// ---------------------------------------------------------------------------
// Extraction
// ---------------------------------------------------------------------------

// Keys returns all keys from m as a slice. The order is non-deterministic.
//
//	maputil.Keys(map[string]int{"a": 1, "b": 2}) // ["a", "b"] in any order
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns all values from m as a slice. The order is non-deterministic.
//
//	maputil.Values(map[string]int{"a": 1, "b": 2}) // [1, 2] in any order
func Values[K comparable, V any](m map[K]V) []V {
	vals := make([]V, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}

// ---------------------------------------------------------------------------
// Merging
// ---------------------------------------------------------------------------

// Merge combines multiple maps into one. Later maps win on key conflict.
// Returns an empty (non-nil) map when called with zero arguments.
//
//	merged := maputil.Merge(defaults, overrides) // overrides win
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// Filtering & Transformation
// ---------------------------------------------------------------------------

// Filter returns a new map containing only key-value pairs for which fn
// returns true.
//
//	positive := maputil.Filter(scores, func(k string, v int) bool { return v > 0 })
func Filter[K comparable, V any](m map[K]V, fn func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if fn(k, v) {
			result[k] = v
		}
	}
	return result
}

// MapValues transforms every value in m using fn, returning a new map with
// the same keys but different value types.
//
//	lengths := maputil.MapValues(names, func(v string) int { return len(v) })
func MapValues[K comparable, V any, U any](m map[K]V, fn func(V) U) map[K]U {
	result := make(map[K]U, len(m))
	for k, v := range m {
		result[k] = fn(v)
	}
	return result
}

// ---------------------------------------------------------------------------
// Inversion & Subsetting
// ---------------------------------------------------------------------------

// Invert swaps keys and values, returning a new map[V]K. If multiple keys map
// to the same value, the last encountered key (non-deterministic) wins.
//
//	maputil.Invert(map[string]int{"a": 1, "b": 2}) // map[int]string{1:"a", 2:"b"}
func Invert[K, V comparable](m map[K]V) map[V]K {
	result := make(map[V]K, len(m))
	for k, v := range m {
		result[v] = k
	}
	return result
}

// Pick returns a new map containing only the specified keys.
//
//	maputil.Pick(user, "name", "email")
func Pick[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	result := make(map[K]V, len(keys))
	for _, k := range keys {
		if v, ok := m[k]; ok {
			result[k] = v
		}
	}
	return result
}

// Omit returns a new map excluding the specified keys.
//
//	maputil.Omit(user, "password", "secret")
func Omit[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	exclude := make(map[K]struct{}, len(keys))
	for _, k := range keys {
		exclude[k] = struct{}{}
	}
	result := make(map[K]V)
	for k, v := range m {
		if _, skip := exclude[k]; !skip {
			result[k] = v
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// Lookup
// ---------------------------------------------------------------------------

// Contains reports whether key exists in m.
//
//	if maputil.Contains(config, "timeout") { ... }
func Contains[K comparable, V any](m map[K]V, key K) bool {
	_, ok := m[key]
	return ok
}

// GetOrDefault returns the value for key, or defaultVal if the key is absent.
// Safe to call on a nil map.
//
//	timeout := maputil.GetOrDefault(config, "timeout", 30)
func GetOrDefault[K comparable, V any](m map[K]V, key K, defaultVal V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return defaultVal
}
