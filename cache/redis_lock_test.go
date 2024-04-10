package cache

import (
	"context"
	"time"
)

func ExampleRefresh() {
	// 假如说我们拿到了一个锁
	var l *Lock

	stop := make(chan struct{}, 1)

	bizStop := make(chan struct{}, 1)
	retryCnt := 0
	go func() {
		// TODO:这里时间一般要根据过期时间来调整，30秒只是举例
		ticker := time.NewTicker(time.Second * 30)
		defer ticker.Stop()
		// 不断续约，直到收到退出信号
		ch := make(chan struct{}, 1)
		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				err := l.Refresh(ctx)
				cancel()
				if err == context.DeadlineExceeded {
					// 一直重试失败怎么办
					ch <- struct{}{}
					continue
				}
				if err != nil {
					// 不可挽回的错误、怎么处理
					bizStop <- struct{}{}
					return
				}
				retryCnt = 0
			case <-ch:
				retryCnt++
				// 重试信号
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				err := l.Refresh(ctx)
				cancel()
				if err == context.DeadlineExceeded {
					// 一直重试失败怎么办
					if retryCnt > 10 {
						// 考虑中断业务
						bizStop <- struct{}{}
						return
					} else {
						ch <- struct{}{}
					}
					continue
				}
				if err != nil {
					// 不可挽回的错误、怎么处理
					//	考虑中断业务
					bizStop <- struct{}{}
					return
				}
				retryCnt = 0
			case <-stop:
				return
			}
		}
	}()

	// 这里是业务代码
	for {
		select {
		case <-bizStop:
			//	回滚操作
		default:
			//	业务代码一步一步循环执行
		}
	}
	// 业务结束 通知不需要再续约
	stop <- struct{}{}
}
