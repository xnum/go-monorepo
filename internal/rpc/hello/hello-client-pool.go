package hello

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"go-monorepo/pkg/pool"
)

// ClientPool defines hello client pool.
type ClientPool struct {
	pool        *pool.GRPCPool
	callTimeout time.Duration
	addr        string
}

// NewClientPoolFromConfig creates hello client.
func NewClientPoolFromConfig(cfg *ClientConfig) *ClientPool {
	q := &ClientPool{
		callTimeout: cfg.CallTimeout,
		addr:        cfg.HelloServiceEndpoint,
	}

	options := pool.NewOptions(cfg.poolConfig,
		[]string{cfg.HelloServiceEndpoint})

	p := pool.NewGRPCPool(options, grpc.WithInsecure())

	q.pool = p

	return q
}

// SayHello calls grpc.
func (c *ClientPool) SayHello(ctx context.Context,
	name string) (string, error) {
	conn, err := c.pool.Get()
	if err != nil {
		return "", fmt.Errorf("dial: %w", err)
	}
	defer c.pool.Put(conn)

	ctx, cancel := context.WithTimeout(ctx, c.callTimeout)
	defer cancel()
	// create client.
	client := NewGreeterClient(conn)

	resp, err := client.SayHello(ctx, &HelloRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	}

	return resp.Message, nil
}
