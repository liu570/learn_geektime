package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"learn_geektime/micro/registry"
	"sync"
)

type Registry struct {
	c       *clientv3.Client
	sess    *concurrency.Session
	cancels []context.CancelFunc
	mutex   sync.Mutex
}

func NewRegistry(c *clientv3.Client) (*Registry, error) {
	// 使用 etcd 自带的租约机制 创建 session
	sess, err := concurrency.NewSession(c)
	if err != nil {
		return nil, err
	}
	return &Registry{
		c:    c,
		sess: sess,
	}, nil
}

func (r *Registry) Register(ctx context.Context, si registry.ServiceInstance) error {
	val, err := json.Marshal(si)
	if err != nil {
		return err
	}
	_, err = r.c.Put(ctx, r.instanceKey(si), string(val), clientv3.WithLease(r.sess.Lease()))
	return err
}

func (r *Registry) UnRegister(ctx context.Context, si registry.ServiceInstance) error {
	_, err := r.c.Delete(ctx, r.instanceKey(si))
	return err
}

func (r *Registry) ListServices(ctx context.Context, serviceName string) ([]registry.ServiceInstance, error) {
	getResp, err := r.c.Get(ctx, r.serviceKey(serviceName), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	res := make([]registry.ServiceInstance, 0, len(getResp.Kvs))
	for _, kv := range getResp.Kvs {
		var si registry.ServiceInstance
		if err := json.Unmarshal(kv.Value, &si); err != nil {
			return nil, err
		}
		res = append(res, si)
	}
	return res, nil

}

func (r *Registry) Subscribe(serviceName string) (<-chan registry.Event, error) {
	ctx, cancel := context.WithCancel(context.Background())
	r.mutex.Lock()
	r.cancels = append(r.cancels, cancel)
	r.mutex.Unlock()
	ctx = clientv3.WithRequireLeader(ctx)
	watchResp := r.c.Watch(ctx, r.serviceKey(serviceName), clientv3.WithPrefix())
	res := make(chan registry.Event)
	go func() {
		for {
			select {
			case resp := <-watchResp:
				if err := resp.Err(); err != nil {
					return
				}
				if resp.Canceled {
					return
				}
				for range resp.Events {
					res <- registry.Event{}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return res, nil
}

func (r *Registry) Close() error {
	r.mutex.Lock()
	cancels := r.cancels
	r.cancels = nil
	r.mutex.Unlock()
	for _, cancel := range cancels {
		cancel()
	}
	err := r.sess.Close()
	return err
}

func (r *Registry) instanceKey(si registry.ServiceInstance) string {
	// 可以直接考虑 Address 字段
	// 也可以考虑在 si 里面引入一个 instanceName 的字段
	return fmt.Sprintf("/micro/%s/%s", si.Name, si.Address)
}

func (r *Registry) serviceKey(serviceName string) string {
	return fmt.Sprintf("/micro/%s", serviceName)
}
