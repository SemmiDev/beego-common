// Package envutil provides helpers for reading environment variables with
// type coercion and default values.
//
// Usage:
//
//	port := envutil.Get("PORT", "8080")
//	timeout := envutil.GetDuration("TIMEOUT", 30*time.Second)
//	debug := envutil.GetBool("DEBUG", false)
//	if err := envutil.Require("DATABASE_URL", "SECRET_KEY"); err != nil {
//	    log.Fatal(err)
//	}
package envutil

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Get returns the value of the environment variable named key, or defaultVal
// when the variable is unset or empty.
//
//	port := envutil.Get("PORT", "8080")
func Get(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// MustGet returns the value of the environment variable named key.
// Panics when the variable is unset or empty.
//
//	dsn := envutil.MustGet("DATABASE_URL")
func MustGet(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("envutil: required environment variable %q is not set", key))
	}
	return v
}

// GetInt returns the value of the environment variable named key as an int.
// Returns defaultVal when the variable is unset, empty, or cannot be parsed.
//
//	workers := envutil.GetInt("WORKER_COUNT", 4)
func GetInt(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return n
}

// GetBool returns the value of the environment variable named key as a bool.
// Recognised truthy values (case-insensitive): "1", "true", "yes", "on".
// Returns defaultVal when the variable is unset or empty.
//
//	debug := envutil.GetBool("DEBUG", false)
func GetBool(key string, defaultVal bool) bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if v == "" {
		return defaultVal
	}
	switch v {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return defaultVal
	}
}

// GetDuration returns the value of the environment variable named key as a
// time.Duration (parsed by time.ParseDuration). Returns defaultVal when the
// variable is unset, empty, or cannot be parsed.
//
//	timeout := envutil.GetDuration("TIMEOUT", 30*time.Second)
func GetDuration(key string, defaultVal time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return defaultVal
	}
	return d
}

// Require checks that every named environment variable is set and non-empty.
// Returns a descriptive error listing all missing variables.
//
//	if err := envutil.Require("DATABASE_URL", "SECRET_KEY"); err != nil {
//	    log.Fatal(err)
//	}
func Require(keys ...string) error {
	var missing []string
	for _, k := range keys {
		if os.Getenv(k) == "" {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("envutil: missing required environment variable(s): %s",
			strings.Join(missing, ", "))
	}
	return nil
}
