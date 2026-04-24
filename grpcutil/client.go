package grpcutil

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// ClientConfig holds dial options for building a gRPC client connection.
type ClientConfig struct {
	// Target address, e.g. "localhost:50051" or "dns:///svc:50051"
	Target string

	// Insecure disables TLS (useful for local dev/testing).
	Insecure bool

	// TLSConfig is used when Insecure is false. If nil, system CAs are used.
	TLSConfig *tls.Config

	// DialTimeout is the max time to wait for a connection to be established.
	// Defaults to 5s.
	DialTimeout time.Duration

	// KeepaliveTime is how often the client pings the server.
	KeepaliveTime time.Duration

	// KeepaliveTimeout is how long the client waits for a keepalive ack.
	KeepaliveTimeout time.Duration

	// MaxRecvMsgSize sets the max message size in bytes the client can receive.
	MaxRecvMsgSize int

	// MaxSendMsgSize sets the max message size in bytes the client can send.
	MaxSendMsgSize int

	// WithBlock makes Dial block until the connection is established.
	WithBlock bool
}

func (c *ClientConfig) defaults() {
	if c.DialTimeout == 0 {
		c.DialTimeout = 5 * time.Second
	}
	if c.KeepaliveTime == 0 {
		c.KeepaliveTime = 30 * time.Second
	}
	if c.KeepaliveTimeout == 0 {
		c.KeepaliveTimeout = 5 * time.Second
	}
	if c.MaxRecvMsgSize == 0 {
		c.MaxRecvMsgSize = 4 * 1024 * 1024
	}
	if c.MaxSendMsgSize == 0 {
		c.MaxSendMsgSize = 4 * 1024 * 1024
	}
}

// Dial creates a new gRPC client connection using the provided config.
func Dial(cfg ClientConfig, extraOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	cfg.defaults()

	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                cfg.KeepaliveTime,
			Timeout:             cfg.KeepaliveTimeout,
			PermitWithoutStream: true,
		}),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(cfg.MaxRecvMsgSize),
			grpc.MaxCallSendMsgSize(cfg.MaxSendMsgSize),
		),
	}

	if cfg.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		tlsCfg := cfg.TLSConfig
		if tlsCfg == nil {
			tlsCfg = &tls.Config{MinVersion: tls.VersionTLS12}
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)))
	}

	if cfg.WithBlock {
		opts = append(opts, grpc.WithBlock())
	}

	opts = append(opts, extraOpts...)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, cfg.Target, opts...) //nolint:staticcheck
	if err != nil {
		return nil, fmt.Errorf("grpcutil: dial %s: %w", cfg.Target, err)
	}
	return conn, nil
}

// MustDial is like Dial but panics on error. Useful in main() setup.
func MustDial(cfg ClientConfig, extraOpts ...grpc.DialOption) *grpc.ClientConn {
	conn, err := Dial(cfg, extraOpts...)
	if err != nil {
		panic(err)
	}
	return conn
}

// DialInsecure is a convenience wrapper for local/dev connections.
func DialInsecure(target string) (*grpc.ClientConn, error) {
	return Dial(ClientConfig{Target: target, Insecure: true})
}
