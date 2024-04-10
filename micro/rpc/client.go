package rpc

import (
	"context"
	"errors"
	"github.com/silenceper/pool"
	"learn_geektime/micro/rpc/message"
	"learn_geektime/micro/rpc/serialize"
	"net"
	"reflect"
	"time"
)

// InitService 为 GetById 之类的字段赋值
func (c *Client) InitService(service Service) error {
	return setFuncField(service, c, c.serializer)
}

func setFuncField(service Service, p Proxy, s serialize.Serializer) error {
	if service == nil {
		return errors.New("rpc:不支持nil")
	}
	// 校验 他是指向结构体的指针
	val := reflect.ValueOf(service)
	typ := reflect.TypeOf(service)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return errors.New("rpc:只支持指向结构体的一级指针")
	}
	val = val.Elem()
	typ = typ.Elem()

	numField := val.NumField()
	for i := 0; i < numField; i++ {
		fieldType := typ.Field(i)
		fieldValue := val.Field(i)

		if !fieldValue.CanSet() {
			continue
		}
		if fieldType.Type.Kind() != reflect.Func {
			continue
		}
		//	替换新的实现
		fn := func(args []reflect.Value) (results []reflect.Value) {
			// 在这里拼凑调用信息，服务名，方法名，参数值

			retVal := reflect.New(fieldType.Type.Out(0).Elem())

			ctx := args[0].Interface().(context.Context)
			reqData, err := s.Encode(args[1].Interface())

			if err != nil {
				return []reflect.Value{retVal, reflect.ValueOf(err)}
			}
			// 获取本地调用信息
			req := &message.Request{
				ServiceName: service.Name(),
				MethodName:  fieldType.Name,
				Data:        reqData,
				Serializer:  s.Code(),
			}
			// 发起远程调用
			resp, err := p.Invoke(ctx, req)
			if err != nil {
				return []reflect.Value{retVal, reflect.ValueOf(err)}
			}
			// 处理远程连接的业务 error
			var er error
			if len(resp.Error) > 0 {
				// 服务端 error
				er = errors.New(string(resp.Error))
			}

			// 这里怎么办
			if len(resp.Data) > 0 {
				err = s.Decode(resp.Data, retVal.Interface())
				if err != nil {
					// 反序列化的 error
					return []reflect.Value{retVal, reflect.ValueOf(err)}
				}
			}
			var retErr reflect.Value
			if er == nil {
				retErr = reflect.Zero(reflect.TypeOf(new(error)).Elem())
			} else {
				retErr = reflect.ValueOf(er)
			}
			return []reflect.Value{retVal, retErr}
		}
		fnVal := reflect.MakeFunc(fieldType.Type, fn)
		fieldValue.Set(fnVal)
	}
	return nil
}

const numOfLengthBytes = 8

type Client struct {
	network string
	addr    string

	pool       pool.Pool
	serializer serialize.Serializer
}

func NewClient(addr string) (*Client, error) {
	p, err := pool.NewChannelPool(&pool.Config{
		InitialCap:  1,
		MaxCap:      30,
		MaxIdle:     10,
		IdleTimeout: time.Minute,
		Factory: func() (interface{}, error) {
			return net.DialTimeout("tcp", addr, time.Second*3)
		},
		Close: func(i interface{}) error {
			return i.(net.Conn).Close()
		},
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		addr:    addr,
		network: "tcp",
		pool:    p,
	}, nil
}

func (c *Client) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	data := message.EncodeReq(req)
	// 新建连接
	respBs, err := c.Send(data)
	if err != nil {
		return nil, err
	}
	return message.DecodeResp(respBs), nil
}

func (c *Client) Send(data []byte) ([]byte, error) {
	val, err := c.pool.Get()
	if err != nil {
		return nil, err
	}
	conn := val.(net.Conn)
	defer func() {
		_ = conn.Close()
	}()
	// 发送请求数据
	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}
	return ReadMsg(conn)
}
