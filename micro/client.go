package micro

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"learn_geektime/micro/registry"
	"time"
)

type ClientOption func(c *Client)

type Client struct {
	insecure bool
	r        registry.Registry
	timeout  time.Duration
	balancer balancer.Builder
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

func ClientWithPickerBuilder(name string, pikcerBuilder base.PickerBuilder) ClientOption {
	return func(c *Client) {
		balanceBuilder := base.NewBalancerBuilder(name, pikcerBuilder, base.Config{HealthCheck: true})
		balancer.Register(balanceBuilder)
		c.balancer = balanceBuilder
	}
}
func (c *Client) Dial(ctx context.Context, service string, dialOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if c.balancer != nil {
		opts = append(opts, grpc.WithDefaultServiceConfig(
			fmt.Sprintf(`{"LoadBalancingPolicy":"%s"}`, c.balancer.Name())))
	}
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
	if dialOpts != nil {
		opts = append(opts, dialOpts...)
	}
	cc, err := grpc.DialContext(ctx, fmt.Sprintf("registry:///%s", service), opts...)
	//cc, err := grpc.DialContext(ctx, "localhost:8081", opts...)
	return cc, err
}
