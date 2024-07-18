package net

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"
)

// 该代码实现了一个简单的连接池

type Pool struct {
	// 空闲连接队列
	idleConns chan *idleConn
	// 请求连接对列
	reqChans chan *connReq

	// 最大连接数量
	maxCnt int
	// 当前连接数量
	cnt int
	// 最大空闲时间
	maxIdleTime time.Duration

	//
	factory func() (net.Conn, error)

	lock sync.Mutex
}

func NewPool(initCnt int, maxCnt int, idleCnt int, maxIdleTime time.Duration, factory func() (net.Conn, error)) (*Pool, error) {
	idleConns := make(chan *idleConn, idleCnt)
	reqChans := make(chan *connReq, maxCnt)
	if initCnt > maxCnt {
		return nil, errors.New("micro:初始连接数量不能大于最大连接数量")
	}

	for i := 0; i < initCnt; i++ {
		conn, err := factory()
		if err != nil {
			return nil, err
		}
		idleConns <- &idleConn{
			conn:           conn,
			lastActiveTime: time.Now(),
		}
	}
	return &Pool{
		idleConns:   idleConns,
		reqChans:    reqChans,
		maxCnt:      maxCnt,
		cnt:         0,
		maxIdleTime: maxIdleTime,
		factory:     factory,
	}, nil
}

func (p *Pool) Get(ctx context.Context) (net.Conn, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	for {
		select {
		case c := <-p.idleConns:
			// 如果空闲队列的连接过期了
			if c.lastActiveTime.Add(p.maxIdleTime).Before(time.Now()) {
				_ = c.conn.Close()
				continue
			}
			return c.conn, nil
		default:
			// 如果没有空闲连接
			p.lock.Lock()
			if p.cnt >= p.maxCnt {
				req := &connReq{connChan: make(chan net.Conn, 1)}
				p.reqChans <- req
				p.lock.Unlock()
				// 下面已经阻塞住了需要释放锁让其它进程可以获得锁
				select {
				// 超时
				case <-ctx.Done():
					// 选项1:删除请求信号
					// 选项2:转发
					go func() {
						c := <-req.connChan
						_ = p.Put(context.Background(), c)
					}()
				//	归还
				case c := <-req.connChan:
					return c, nil
				}

			} else {
				conn, err := p.factory()
				if err != nil {
					return nil, err
				}
				p.cnt++
				p.lock.Unlock()
				return conn, err
			}
		}
	}
}

func (p *Pool) Put(ctx context.Context, conn net.Conn) error {
	if len(p.reqChans) >= 0 {
		req := <-p.reqChans
		req.connChan <- conn
		return nil
	}
	ic := &idleConn{
		conn:           conn,
		lastActiveTime: time.Now(),
	}

	select {
	case p.idleConns <- ic:
	default:
		conn.Close()
		p.lock.Lock()
		p.cnt--
		p.lock.Unlock()
	}
	return nil
}

type idleConn struct {
	conn           net.Conn
	lastActiveTime time.Time
}

type connReq struct {
	connChan chan net.Conn
}
