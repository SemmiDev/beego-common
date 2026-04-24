package grpcutil

import (
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ---- Unary Interceptors ----

// UnaryLoggerInterceptor logs each unary RPC call with duration and status.
func UnaryLoggerInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		code, msg := FromError(err)
		logArgs := []interface{}{
			"method", info.FullMethod,
			"duration_ms", duration.Milliseconds(),
			"code", code.String(),
		}
		if err != nil {
			logArgs = append(logArgs, "error", msg)
			logger.ErrorContext(ctx, "grpc unary call failed", logArgs...)
		} else {
			logger.InfoContext(ctx, "grpc unary call", logArgs...)
		}
		return resp, err
	}
}

// UnaryRecoveryInterceptor recovers from panics in unary handlers and returns
// an Internal gRPC error instead of crashing the server.
func UnaryRecoveryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				if logger != nil {
					logger.ErrorContext(ctx, "panic recovered",
						"method", info.FullMethod,
						"panic", r,
						"stack", string(stack),
					)
				}
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}

// UnaryTimeoutInterceptor wraps each unary RPC with a deadline if none is set.
func UnaryTimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		return handler(ctx, req)
	}
}

// UnaryAuthInterceptor validates incoming requests using a custom AuthFunc.
// AuthFunc receives the context and full method name; return error to reject.
type AuthFunc func(ctx context.Context, fullMethod string) (context.Context, error)

func UnaryAuthInterceptor(authFn AuthFunc) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		newCtx, err := authFn(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// ---- Stream Interceptors ----

// StreamLoggerInterceptor logs each streaming RPC call.
func StreamLoggerInterceptor(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()
		err := handler(srv, ss)
		duration := time.Since(start)

		code, msg := FromError(err)
		logArgs := []interface{}{
			"method", info.FullMethod,
			"duration_ms", duration.Milliseconds(),
			"code", code.String(),
			"client_stream", info.IsClientStream,
			"server_stream", info.IsServerStream,
		}
		if err != nil {
			logArgs = append(logArgs, "error", msg)
			logger.ErrorContext(ss.Context(), "grpc stream call failed", logArgs...)
		} else {
			logger.InfoContext(ss.Context(), "grpc stream call", logArgs...)
		}
		return err
	}
}

// StreamRecoveryInterceptor recovers from panics in stream handlers.
func StreamRecoveryInterceptor(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				if logger != nil {
					logger.ErrorContext(ss.Context(), "panic recovered in stream",
						"method", info.FullMethod,
						"panic", r,
						"stack", string(stack),
					)
				}
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		return handler(srv, ss)
	}
}

// StreamAuthInterceptor validates streaming RPCs using a custom AuthFunc.
func StreamAuthInterceptor(authFn AuthFunc) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		newCtx, err := authFn(ss.Context(), info.FullMethod)
		if err != nil {
			return err
		}
		return handler(srv, &wrappedStream{ss, newCtx})
	}
}

// wrappedStream injects a new context into a grpc.ServerStream.
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context { return w.ctx }

// ---- Chain helpers ----

// ChainUnaryInterceptors chains multiple unary interceptors into one.
// Execution order: first → last (outermost → innermost).
//
// Note: grpc.ChainUnaryInterceptor returns a grpc.ServerOption, not a
// grpc.UnaryServerInterceptor, so we implement chaining manually here.
func ChainUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	n := len(interceptors)
	switch n {
	case 0:
		return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			return handler(ctx, req)
		}
	case 1:
		return interceptors[0]
	default:
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			// Build the chain from innermost (last) to outermost (first).
			chained := handler
			for i := n - 1; i > 0; i-- {
				idx := i
				prev := chained
				chained = func(ctx context.Context, req interface{}) (interface{}, error) {
					return interceptors[idx](ctx, req, info, prev)
				}
			}
			return interceptors[0](ctx, req, info, chained)
		}
	}
}

// ChainStreamInterceptors chains multiple stream interceptors into one.
// Execution order: first → last (outermost → innermost).
//
// Note: grpc.ChainStreamInterceptor returns a grpc.ServerOption, not a
// grpc.StreamServerInterceptor, so we implement chaining manually here.
func ChainStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	n := len(interceptors)
	switch n {
	case 0:
		return func(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			return handler(srv, ss)
		}
	case 1:
		return interceptors[0]
	default:
		return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			// Build the chain from innermost (last) to outermost (first).
			chained := handler
			for i := n - 1; i > 0; i-- {
				idx := i
				prev := chained
				chained = func(srv interface{}, ss grpc.ServerStream) error {
					return interceptors[idx](srv, ss, info, prev)
				}
			}
			return interceptors[0](srv, ss, info, chained)
		}
	}
}
