package grpcutil

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ---- Validation helpers ----

// ValidationError accumulates field-level validation errors and returns
// a single gRPC InvalidArgument status with a descriptive message.
type ValidationError struct {
	errs []string
}

// Field records a validation failure for the given field.
func (v *ValidationError) Field(field, reason string) *ValidationError {
	v.errs = append(v.errs, field+": "+reason)
	return v
}

// Required records an error if val is empty.
func (v *ValidationError) Required(field, val string) *ValidationError {
	if strings.TrimSpace(val) == "" {
		v.Field(field, "is required")
	}
	return v
}

// MinLen records an error if val is shorter than min characters.
func (v *ValidationError) MinLen(field, val string, min int) *ValidationError {
	if len(strings.TrimSpace(val)) < min {
		v.Field(field, "is too short")
	}
	return v
}

// MaxLen records an error if val is longer than max characters.
func (v *ValidationError) MaxLen(field, val string, max int) *ValidationError {
	if len(val) > max {
		v.Field(field, "is too long")
	}
	return v
}

// Positive records an error if val is not > 0.
func (v *ValidationError) Positive(field string, val int64) *ValidationError {
	if val <= 0 {
		v.Field(field, "must be positive")
	}
	return v
}

// HasErrors reports whether any validation errors were recorded.
func (v *ValidationError) HasErrors() bool {
	return len(v.errs) > 0
}

// Err returns a gRPC InvalidArgument error if any validation errors exist,
// otherwise returns nil.
func (v *ValidationError) Err() error {
	if !v.HasErrors() {
		return nil
	}
	msg := "validation failed: " + strings.Join(v.errs, "; ")
	return status.Error(codes.InvalidArgument, msg)
}

// Validate is a convenience constructor — create, populate, and return error.
//
//	if err := grpcutil.Validate().Required("name", req.Name).Err(); err != nil {
//	    return nil, err
//	}
func Validate() *ValidationError {
	return &ValidationError{}
}

// ---- Context helpers ----

// ContextError maps context errors to appropriate gRPC status errors.
func ContextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return status.Error(codes.Canceled, "request canceled")
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, "deadline exceeded")
	default:
		return nil
	}
}

// IsCanceled reports whether a gRPC call was cancelled by the client.
func IsCanceled(err error) bool {
	return IsCode(err, codes.Canceled)
}

// IsDeadlineExceeded reports whether a gRPC call exceeded its deadline.
func IsDeadlineExceeded(err error) bool {
	return IsCode(err, codes.DeadlineExceeded)
}

// ---- Retry helpers ----

// IsRetryable reports whether the gRPC error is considered safe to retry.
// Retryable codes: Unavailable, ResourceExhausted, DeadlineExceeded.
func IsRetryable(err error) bool {
	code, _ := FromError(err)
	switch code {
	case codes.Unavailable, codes.ResourceExhausted, codes.DeadlineExceeded:
		return true
	default:
		return false
	}
}
