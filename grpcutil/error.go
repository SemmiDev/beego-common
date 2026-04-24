package grpcutil

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCError wraps a gRPC status error with additional context.
type GRPCError struct {
	Code    codes.Code
	Message string
	Details []interface{}
}

func (e *GRPCError) Error() string {
	return e.Message
}

func (e *GRPCError) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Message)
}

// NewError creates a gRPC status error.
func NewError(code codes.Code, msg string) error {
	return status.Error(code, msg)
}

// NewErrorf creates a gRPC status error with formatting.
func NewErrorf(code codes.Code, format string, args ...interface{}) error {
	return status.Errorf(code, format, args...)
}

// Common gRPC error constructors — mirrors common HTTP error helpers.

func ErrNotFound(msg string) error        { return status.Error(codes.NotFound, msg) }
func ErrInvalidArg(msg string) error      { return status.Error(codes.InvalidArgument, msg) }
func ErrAlreadyExists(msg string) error   { return status.Error(codes.AlreadyExists, msg) }
func ErrUnauthenticated(msg string) error { return status.Error(codes.Unauthenticated, msg) }
func ErrPermDenied(msg string) error      { return status.Error(codes.PermissionDenied, msg) }
func ErrInternal(msg string) error        { return status.Error(codes.Internal, msg) }
func ErrUnavailable(msg string) error     { return status.Error(codes.Unavailable, msg) }
func ErrUnimplemented(msg string) error   { return status.Error(codes.Unimplemented, msg) }
func ErrDeadlineExceeded(msg string) error {
	return status.Error(codes.DeadlineExceeded, msg)
}
func ErrResourceExhausted(msg string) error {
	return status.Error(codes.ResourceExhausted, msg)
}

// FromError extracts gRPC status code and message from an error.
// Returns (codes.OK, "") if err is nil.
func FromError(err error) (codes.Code, string) {
	if err == nil {
		return codes.OK, ""
	}
	if s, ok := status.FromError(err); ok {
		return s.Code(), s.Message()
	}
	return codes.Internal, err.Error()
}

// IsCode reports whether err is a gRPC status error with the given code.
func IsCode(err error, code codes.Code) bool {
	s, ok := status.FromError(err)
	return ok && s.Code() == code
}

// IsNotFound reports whether err is a gRPC NotFound error.
func IsNotFound(err error) bool { return IsCode(err, codes.NotFound) }

// IsInvalidArg reports whether err is a gRPC InvalidArgument error.
func IsInvalidArg(err error) bool { return IsCode(err, codes.InvalidArgument) }

// IsAlreadyExists reports whether err is a gRPC AlreadyExists error.
func IsAlreadyExists(err error) bool { return IsCode(err, codes.AlreadyExists) }

// IsUnauthenticated reports whether err is a gRPC Unauthenticated error.
func IsUnauthenticated(err error) bool { return IsCode(err, codes.Unauthenticated) }

// IsInternal reports whether err is a gRPC Internal error.
func IsInternal(err error) bool { return IsCode(err, codes.Internal) }

// HTTPStatusToCode maps HTTP status codes to gRPC codes.
func HTTPStatusToCode(httpStatus int) codes.Code {
	switch httpStatus {
	case http.StatusOK:
		return codes.OK
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded
	case http.StatusNotImplemented:
		return codes.Unimplemented
	default:
		return codes.Unknown
	}
}

// CodeToHTTPStatus maps gRPC codes to HTTP status codes.
func CodeToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.InvalidArgument, codes.FailedPrecondition, codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists, codes.Aborted:
		return http.StatusConflict
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Canceled:
		return 499 // Client Closed Request
	case codes.DataLoss, codes.Unknown, codes.Internal:
		return http.StatusInternalServerError
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}

// WrapError wraps a standard error into an appropriate gRPC status error.
// Useful when bridging domain/apperror types to gRPC.
func WrapError(err error, fallbackCode codes.Code) error {
	if err == nil {
		return nil
	}
	// If already a gRPC status error, pass through.
	if _, ok := status.FromError(err); ok {
		return err
	}
	// Map sentinel errors if needed.
	var target *GRPCError
	if errors.As(err, &target) {
		return target.GRPCStatus().Err()
	}
	return status.Error(fallbackCode, err.Error())
}
