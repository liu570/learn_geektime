package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/silenceper/pool"
	"net"
	"reflect"
	"time"
)

// InitClientProxy 为 GetById 之类的字段赋值
func InitClientProxy(addr string, service Service) error {
	client, err := NewClient(addr)
	if err != nil {
		return err
	}
	return setFuncField(service, client)
}

func setFuncField(service Service, p Proxy) error {
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
			reqData, err := json.Marshal(args[1].Interface())
			if err != nil {
				return []reflect.Value{retVal, reflect.ValueOf(err)}
			}
			// 获取本地调用信息
			req := &Request{
				ServiceName: service.Name(),
				MethodName:  fieldType.Name,
				Arg:         reqData,
			}
			// 发起远程调用
			resp, err := p.Invoke(ctx, req)
			if err != nil {
				return []reflect.Value{retVal, reflect.ValueOf(err)}
			}
			// 这里怎么办
			fmt.Println("client-resp:", resp)
			err = json.Unmarshal(resp.data, retVal.Interface())
			if err != nil {
				return []reflect.Value{retVal, reflect.ValueOf(err)}
			}
			return []reflect.Value{retVal, reflect.Zero(reflect.TypeOf(new(error)).Elem())}
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

	pool pool.Pool
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

func (c *Client) Invoke(ctx context.Context, req *Request) (*Response, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// 新建连接
	resp, err := c.Send(data)
	if err != nil {
		return nil, err
	}
	return &Response{
		data: resp,
	}, nil
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
	req := EncodeMsg(data)
	// 发送请求数据
	_, err = conn.Write(req)
	if err != nil {
		return nil, err
	}
	return ReadMsg(conn)
}
