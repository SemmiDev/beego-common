// Package sliceutil provides generic utility functions for working with
// slices. All functions use Go 1.18+ generics.
//
// Usage:
//
//	ids := sliceutil.Map(users, func(u User) string { return u.ID })
//	active := sliceutil.Filter(users, func(u User) bool { return u.Active })
//	if sliceutil.Contains(roles, "admin") { ... }
//	chunks := sliceutil.Chunk(items, 100) // batch processing
package sliceutil

// ---------------------------------------------------------------------------
// Transformation
// ---------------------------------------------------------------------------

// Map applies fn to each element of src and returns a new slice of results.
//
//	names := sliceutil.Map(users, func(u User) string { return u.Name })
func Map[T any, U any](src []T, fn func(T) U) []U {
	if src == nil {
		return nil
	}
	result := make([]U, len(src))
	for i, v := range src {
		result[i] = fn(v)
	}
	return result
}

// Filter returns a new slice containing only elements for which fn returns true.
//
//	admins := sliceutil.Filter(users, func(u User) bool { return u.Role == "admin" })
func Filter[T any](src []T, fn func(T) bool) []T {
	if src == nil {
		return nil
	}
	result := make([]T, 0, len(src)/2)
	for _, v := range src {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce reduces a slice to a single value by applying fn to an accumulator
// and each element.
//
//	total := sliceutil.Reduce(prices, 0, func(acc int, p Price) int { return acc + p.Amount })
func Reduce[T any, U any](src []T, initial U, fn func(U, T) U) U {
	acc := initial
	for _, v := range src {
		acc = fn(acc, v)
	}
	return acc
}

// FlatMap applies fn to each element and flattens the results into a single slice.
func FlatMap[T any, U any](src []T, fn func(T) []U) []U {
	var result []U
	for _, v := range src {
		result = append(result, fn(v)...)
	}
	return result
}

// ---------------------------------------------------------------------------
// Lookup
// ---------------------------------------------------------------------------

// Contains reports whether v is present in the slice.
//
//	if sliceutil.Contains(roles, "admin") { ... }
func Contains[T comparable](src []T, v T) bool {
	for _, item := range src {
		if item == v {
			return true
		}
	}
	return false
}

// IndexOf returns the first index of v in src, or -1 if not found.
func IndexOf[T comparable](src []T, v T) int {
	for i, item := range src {
		if item == v {
			return i
		}
	}
	return -1
}

// Find returns the first element for which fn returns true, along with a bool
// indicating whether an element was found.
func Find[T any](src []T, fn func(T) bool) (T, bool) {
	for _, v := range src {
		if fn(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// ---------------------------------------------------------------------------
// Grouping & Partitioning
// ---------------------------------------------------------------------------

// GroupBy groups elements by a key extracted via fn.
//
//	byDept := sliceutil.GroupBy(employees, func(e Employee) string { return e.Department })
func GroupBy[T any, K comparable](src []T, fn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, v := range src {
		key := fn(v)
		result[key] = append(result[key], v)
	}
	return result
}

// Chunk splits a slice into chunks of the given size. The last chunk may be
// smaller than size. Useful for batch processing.
//
//	for _, batch := range sliceutil.Chunk(items, 100) { processBatch(batch) }
func Chunk[T any](src []T, size int) [][]T {
	if size <= 0 || len(src) == 0 {
		return nil
	}
	chunks := make([][]T, 0, (len(src)+size-1)/size)
	for i := 0; i < len(src); i += size {
		end := i + size
		if end > len(src) {
			end = len(src)
		}
		chunks = append(chunks, src[i:end])
	}
	return chunks
}

// ---------------------------------------------------------------------------
// Deduplication & Set operations
// ---------------------------------------------------------------------------

// Unique returns a new slice with duplicate elements removed, preserving
// the order of first occurrence.
//
//	uniq := sliceutil.Unique([]int{1, 2, 2, 3, 1}) // [1, 2, 3]
func Unique[T comparable](src []T) []T {
	if src == nil {
		return nil
	}
	seen := make(map[T]struct{}, len(src))
	result := make([]T, 0, len(src))
	for _, v := range src {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// Difference returns elements in a that are not in b.
//
//	diff := sliceutil.Difference([]int{1,2,3,4}, []int{2,4}) // [1, 3]
func Difference[T comparable](a, b []T) []T {
	bSet := make(map[T]struct{}, len(b))
	for _, v := range b {
		bSet[v] = struct{}{}
	}
	var result []T
	for _, v := range a {
		if _, ok := bSet[v]; !ok {
			result = append(result, v)
		}
	}
	return result
}

// Intersect returns elements that are present in both a and b.
func Intersect[T comparable](a, b []T) []T {
	bSet := make(map[T]struct{}, len(b))
	for _, v := range b {
		bSet[v] = struct{}{}
	}
	var result []T
	for _, v := range a {
		if _, ok := bSet[v]; ok {
			result = append(result, v)
		}
	}
	return Unique(result)
}

// ---------------------------------------------------------------------------
// Convenience
// ---------------------------------------------------------------------------

// ToMap converts a slice to a map using fn to extract the key for each element.
//
//	byID := sliceutil.ToMap(users, func(u User) string { return u.ID })
func ToMap[T any, K comparable](src []T, fn func(T) K) map[K]T {
	result := make(map[K]T, len(src))
	for _, v := range src {
		result[fn(v)] = v
	}
	return result
}

// Keys extracts all keys from a map as a slice.
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values extracts all values from a map as a slice.
func Values[K comparable, V any](m map[K]V) []V {
	vals := make([]V, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}

// IsEmpty reports whether the slice is nil or has zero length.
func IsEmpty[T any](src []T) bool {
	return len(src) == 0
}

// ---------------------------------------------------------------------------
// Access
// ---------------------------------------------------------------------------

// First returns the first element of src and true.
// Returns the zero value and false when src is nil or empty.
//
//	v, ok := sliceutil.First(items)
func First[T any](src []T) (T, bool) {
	if len(src) == 0 {
		var zero T
		return zero, false
	}
	return src[0], true
}

// Last returns the last element of src and true.
// Returns the zero value and false when src is nil or empty.
//
//	v, ok := sliceutil.Last(items)
func Last[T any](src []T) (T, bool) {
	if len(src) == 0 {
		var zero T
		return zero, false
	}
	return src[len(src)-1], true
}

// ---------------------------------------------------------------------------
// Predicates
// ---------------------------------------------------------------------------

// Any reports whether at least one element satisfies fn.
// Returns false for a nil or empty slice.
//
//	sliceutil.Any(items, func(i Item) bool { return i.Active })
func Any[T any](src []T, fn func(T) bool) bool {
	for _, v := range src {
		if fn(v) {
			return true
		}
	}
	return false
}

// All reports whether every element satisfies fn.
// Returns true vacuously for a nil or empty slice.
//
//	sliceutil.All(items, func(i Item) bool { return i.Active })
func All[T any](src []T, fn func(T) bool) bool {
	for _, v := range src {
		if !fn(v) {
			return false
		}
	}
	return true
}

// None reports whether no element satisfies fn.
// Returns true vacuously for a nil or empty slice.
//
//	sliceutil.None(items, func(i Item) bool { return i.Deleted })
func None[T any](src []T, fn func(T) bool) bool {
	for _, v := range src {
		if fn(v) {
			return false
		}
	}
	return true
}

// Count returns the number of elements that satisfy fn.
//
//	n := sliceutil.Count(items, func(i Item) bool { return i.Active })
func Count[T any](src []T, fn func(T) bool) int {
	n := 0
	for _, v := range src {
		if fn(v) {
			n++
		}
	}
	return n
}

// ---------------------------------------------------------------------------
// Reshaping
// ---------------------------------------------------------------------------

// Reverse returns a reversed copy of src. The original slice is never mutated.
//
//	sliceutil.Reverse([]int{1, 2, 3}) // [3, 2, 1]
func Reverse[T any](src []T) []T {
	if src == nil {
		return nil
	}
	result := make([]T, len(src))
	for i, v := range src {
		result[len(src)-1-i] = v
	}
	return result
}

// Flatten collapses a [][]T into a []T.
//
//	sliceutil.Flatten([][]int{{1, 2}, {3}}) // [1, 2, 3]
func Flatten[T any](src [][]T) []T {
	var result []T
	for _, inner := range src {
		result = append(result, inner...)
	}
	return result
}

// Compact removes zero-value elements from src using generic equality.
//
//	sliceutil.Compact([]string{"a", "", "b", ""}) // ["a", "b"]
func Compact[T comparable](src []T) []T {
	var zero T
	var result []T
	for _, v := range src {
		if v != zero {
			result = append(result, v)
		}
	}
	return result
}

// Take returns the first n elements of src.
// Safe when n >= len(src) — returns a copy of the whole slice.
//
//	sliceutil.Take([]int{1, 2, 3, 4, 5}, 3) // [1, 2, 3]
func Take[T any](src []T, n int) []T {
	if n <= 0 || len(src) == 0 {
		return nil
	}
	if n >= len(src) {
		result := make([]T, len(src))
		copy(result, src)
		return result
	}
	result := make([]T, n)
	copy(result, src[:n])
	return result
}

// Skip drops the first n elements and returns the rest.
// Safe when n >= len(src) — returns nil.
//
//	sliceutil.Skip([]int{1, 2, 3, 4, 5}, 2) // [3, 4, 5]
func Skip[T any](src []T, n int) []T {
	if n <= 0 {
		result := make([]T, len(src))
		copy(result, src)
		return result
	}
	if n >= len(src) {
		return nil
	}
	result := make([]T, len(src)-n)
	copy(result, src[n:])
	return result
}
