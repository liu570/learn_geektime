package session

import (
	"learn_geektime/web"
)

// 做为用户来说 用户一点也不关心我们是如何实现 session 的
// 所以我们可以提供一个工具供用户简洁使用

type Manager struct {
	Store
	Propagator
	SessCtxKey string
}

func (m *Manager) InitSession(ctx *web.Context, id string) (Session, error) {
	sess, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	if err = m.Inject(id, ctx.Resp); err != nil {
		return nil, err
	}
	return sess, nil
}

func (m *Manager) GetSession(ctx *web.Context) (Session, error) {
	sess, ok := ctx.UserValues[m.SessCtxKey]
	if ok {
		return sess.(Session), nil
	}
	id, err := m.Propagator.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}
	sess, err = m.Store.Get(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	ctx.UserValues[m.SessCtxKey] = sess

	return sess.(Session), nil
}

func (m *Manager) RemoveSession(ctx *web.Context) error {
	id, err := m.Propagator.Extract(ctx.Req)
	if err != nil {
		return err
	}
	err = m.Store.Remove(ctx.Req.Context(), id)
	if err != nil {
		return err
	}
	err = m.Propagator.Remove(ctx.Resp)
	if err != nil {
		return err
	}
	delete(ctx.UserValues, m.SessCtxKey)
	return err
}
