package session

import (
	"context"
	"net/http"
)

type Session interface {
	Get(ctx context.Context, key string) (string, error)
	// val 如果设计类型为 any 的话，那么对应的 redis 之类的实现需要考虑序列化的问题
	Set(ctx context.Context, key string, val string) error
	ID() string
}

// Store 大致意思为管理 Session 的人
type Store interface {
	Generate(ctx context.Context, id string) (Session, error)
	Remove(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (Session, error)

	Refresh(ctx context.Context, id string) error
}

// Propagator 在开发过程中 session id 可以存储在许多地方，主流都存储在cookie中
// 所以我们使用 Propagator 做为一个抽象层，不同的实现允许将 session id 存储在不同的位置

type Propagator interface {
	// Inject 注入 session id
	Inject(id string, resp http.ResponseWriter) error

	// Extract 提取 session id
	Extract(req *http.Request) (string, error)

	Remove(writer http.ResponseWriter) error
}
