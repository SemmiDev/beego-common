package grpcutil

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// ServerConfig holds configuration for building a gRPC server.
type ServerConfig struct {
	// Port to listen on, e.g. ":50051"
	Port string

	// MaxRecvMsgSize sets the max message size in bytes the server can receive.
	// Defaults to 4MB.
	MaxRecvMsgSize int

	// MaxSendMsgSize sets the max message size in bytes the server can send.
	// Defaults to 4MB.
	MaxSendMsgSize int

	// KeepaliveTime is how often the server pings idle clients.
	KeepaliveTime time.Duration

	// KeepaliveTimeout is how long the server waits for a keepalive ack.
	KeepaliveTimeout time.Duration

	// EnableReflection enables gRPC server reflection (useful for grpcurl, Postman).
	EnableReflection bool

	// EnableHealthCheck registers the standard gRPC health check service.
	EnableHealthCheck bool

	// Logger for server lifecycle events. Uses slog.Default() if nil.
	Logger *slog.Logger
}

func (c *ServerConfig) defaults() {
	if c.Port == "" {
		c.Port = ":50051"
	}
	if c.MaxRecvMsgSize == 0 {
		c.MaxRecvMsgSize = 4 * 1024 * 1024 // 4MB
	}
	if c.MaxSendMsgSize == 0 {
		c.MaxSendMsgSize = 4 * 1024 * 1024 // 4MB
	}
	if c.KeepaliveTime == 0 {
		c.KeepaliveTime = 30 * time.Second
	}
	if c.KeepaliveTimeout == 0 {
		c.KeepaliveTimeout = 5 * time.Second
	}
	if c.Logger == nil {
		c.Logger = slog.Default()
	}
}

// ServerBuilder fluently constructs a *grpc.Server.
type ServerBuilder struct {
	cfg          ServerConfig
	unaryInts    []grpc.UnaryServerInterceptor
	streamInts   []grpc.StreamServerInterceptor
	serverOpts   []grpc.ServerOption
	healthServer *health.Server
}

// NewServerBuilder creates a new ServerBuilder with the given config.
func NewServerBuilder(cfg ServerConfig) *ServerBuilder {
	cfg.defaults()
	return &ServerBuilder{cfg: cfg}
}

// WithUnaryInterceptor appends a unary interceptor.
func (b *ServerBuilder) WithUnaryInterceptor(i grpc.UnaryServerInterceptor) *ServerBuilder {
	b.unaryInts = append(b.unaryInts, i)
	return b
}

// WithStreamInterceptor appends a stream interceptor.
func (b *ServerBuilder) WithStreamInterceptor(i grpc.StreamServerInterceptor) *ServerBuilder {
	b.streamInts = append(b.streamInts, i)
	return b
}

// WithServerOption appends a raw grpc.ServerOption.
func (b *ServerBuilder) WithServerOption(opt grpc.ServerOption) *ServerBuilder {
	b.serverOpts = append(b.serverOpts, opt)
	return b
}

// Build constructs the *grpc.Server.
func (b *ServerBuilder) Build() *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(b.cfg.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(b.cfg.MaxSendMsgSize),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    b.cfg.KeepaliveTime,
			Timeout: b.cfg.KeepaliveTimeout,
		}),
	}
	if len(b.unaryInts) > 0 {
		opts = append(opts, grpc.ChainUnaryInterceptor(b.unaryInts...))
	}
	if len(b.streamInts) > 0 {
		opts = append(opts, grpc.ChainStreamInterceptor(b.streamInts...))
	}
	opts = append(opts, b.serverOpts...)

	srv := grpc.NewServer(opts...)

	if b.cfg.EnableReflection {
		reflection.Register(srv)
	}
	if b.cfg.EnableHealthCheck {
		b.healthServer = health.NewServer()
		grpc_health_v1.RegisterHealthServer(srv, b.healthServer)
	}

	return srv
}

// HealthServer returns the registered health server, if any.
func (b *ServerBuilder) HealthServer() *health.Server {
	return b.healthServer
}

// ---- Lifecycle helpers ----

// ListenAndServe starts the gRPC server and blocks until an OS signal
// (SIGINT / SIGTERM) is received, then performs a graceful shutdown.
func ListenAndServe(srv *grpc.Server, cfg ServerConfig) error {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		return fmt.Errorf("grpcutil: listen on %s: %w", cfg.Port, err)
	}

	errCh := make(chan error, 1)
	go func() {
		cfg.Logger.Info("grpc server starting", "port", cfg.Port)
		if err := srv.Serve(lis); err != nil {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		cfg.Logger.Info("grpc server shutting down", "signal", sig)
		srv.GracefulStop()
		return nil
	case err := <-errCh:
		return fmt.Errorf("grpcutil: serve: %w", err)
	}
}

// ListenAndServeContext starts the gRPC server and stops when ctx is cancelled
// or an OS signal is received.
func ListenAndServeContext(ctx context.Context, srv *grpc.Server, cfg ServerConfig) error {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		return fmt.Errorf("grpcutil: listen on %s: %w", cfg.Port, err)
	}

	errCh := make(chan error, 1)
	go func() {
		cfg.Logger.Info("grpc server starting", "port", cfg.Port)
		if err := srv.Serve(lis); err != nil {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		cfg.Logger.Info("grpc server stopping (context cancelled)")
		srv.GracefulStop()
		return ctx.Err()
	case sig := <-quit:
		cfg.Logger.Info("grpc server stopping", "signal", sig)
		srv.GracefulStop()
		return nil
	case err := <-errCh:
		return fmt.Errorf("grpcutil: serve: %w", err)
	}
}
