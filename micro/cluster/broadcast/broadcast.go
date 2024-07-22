package broadcast

import (
	"context"
	"google.golang.org/grpc"
	"learn_geektime/micro/registry"
	"sync"
)

type ClusterBuilder struct {
	registry registry.Registry
	service  string
}

func NewClusterBuilder(registry registry.Registry, serviceName string) *ClusterBuilder {
	return &ClusterBuilder{
		registry: registry,
		service:  serviceName,
	}
}

func (c ClusterBuilder) buildUnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if !isBoardCastKey(ctx) {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		instances, err := c.registry.ListServices(ctx, c.service)
		if err != nil {
			return err
		}
		var wg sync.WaitGroup
		for _, instance := range instances {
			addr := instance.Address
			go func() error {
				wg.Add(1)
				insCc, err := grpc.Dial(addr)
				if err != nil {
					return err
				}
				defer wg.Done()
				// 重新选择连接并进行调用
				return invoker(ctx, method, req, reply, insCc, opts...)
			}()
		}
		wg.Wait()
		return nil
	}
}

func WithBoardCast(ctx context.Context) context.Context {
	return context.WithValue(ctx, boardcastKey{}, true)
}

type boardcastKey struct{}

func isBoardCastKey(ctx context.Context) bool {
	val, ok := ctx.Value(boardcastKey{}).(bool)
	return ok && val
}
