package micro

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"learn_geektime/micro/registry"
	"time"
)

type ClientOption func(c *Client)

type Client struct {
	insecure bool
	r        registry.Registry
	timeout  time.Duration
}

func NewClient(opts ...ClientOption) (*Client, error) {
	client := &Client{}
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}

func ClientWithInsecure() ClientOption {
	return func(c *Client) {
		c.insecure = true
	}
}

func ClientWithRegistry(r registry.Registry, timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.r = r
		c.timeout = timeout
	}
}
func (c *Client) Dial(ctx context.Context, service string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if c.r != nil {
		rb, err := NewResolverBuilder(c.r, c.timeout)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithResolvers(rb))
	}
	if c.insecure {
		opts = append(opts, grpc.WithInsecure())
	}
	cc, err := grpc.DialContext(ctx, fmt.Sprintf("registry:///%s", service), opts...)
	//cc, err := grpc.DialContext(ctx, "localhost:8081", opts...)
	return cc, err
}
