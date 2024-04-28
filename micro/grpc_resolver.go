package micro

import (
	"context"
	"fmt"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"learn_geektime/micro/registry"
	"time"
)

type GrpcResolverBuilder struct {
	r       registry.Registry
	timeout time.Duration
}

func NewResolverBuilder(r registry.Registry, timeout time.Duration) (*GrpcResolverBuilder, error) {
	return &GrpcResolverBuilder{
		r:       r,
		timeout: timeout,
	}, nil
}

func (b *GrpcResolverBuilder) Build(target resolver.Target,
	cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	g := &GrpcResolver{
		cc:      cc,
		r:       b.r,
		target:  target,
		timeout: b.timeout,
	}
	g.resolve()
	go g.watch()
	return g, nil

}

func (b *GrpcResolverBuilder) Scheme() string {
	return "registry"
}

type GrpcResolver struct {
	// "registry:///localhost:8081"
	r       registry.Registry
	cc      resolver.ClientConn
	target  resolver.Target
	timeout time.Duration
	close   chan struct{}
}

func (g *GrpcResolver) ResolveNow(options resolver.ResolveNowOptions) {
	g.resolve()
}

func (g *GrpcResolver) watch() {
	event, err := g.r.Subscribe(g.target.Endpoint())
	if err != nil {
		g.cc.ReportError(err)
		return
	}
	select {
	case <-event:
		g.resolve()
	case <-g.close:
		return
	}
}

func (g *GrpcResolver) resolve() {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	instances, err := g.r.ListServices(ctx, g.target.Endpoint())
	if err != nil {
		g.cc.ReportError(fmt.Errorf("micro:list services failed: %v", err))
		return
	}
	address := make([]resolver.Address, 0, len(instances))
	for _, si := range instances {
		address = append(address, resolver.Address{
			Addr:       si.Address,
			Attributes: attributes.New("weight", si.Weight),
		})
	}
	err = g.cc.UpdateState(resolver.State{Addresses: address})
	if err != nil {
		g.cc.ReportError(err)
		return
	}
}

func (g *GrpcResolver) Close() {
	// 利用关闭 chan 会发送一次信息的特性使得 close chan 只接收一次信息
	close(g.close)
}
