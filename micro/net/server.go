package net

import (
	"encoding/binary"
	"net"
)

const numOfLengthBytes = 8

type Server struct {
	network string
	addr    string
}

func NewServer(network string, addr string) *Server {
	return &Server{
		network: network,
		addr:    addr,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen(s.network, s.addr)
	if err != nil {
		// 常见错误为端口被占用
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			if er := s.handleConn(conn); er != nil {
				_ = conn.Close()
			}
		}()
	}
}

// 我们可以任务 一个请求包含两个部分
// 1. 长度字段：8字节
// 2. 请求数据
func (s Server) handleConn(conn net.Conn) error {
	for {
		// lenBs 长度字段的字节表示
		lenBs := make([]byte, numOfLengthBytes)
		_, err := conn.Read(lenBs)
		if err != nil {
			return err
		}

		// 消息有多长
		length := binary.BigEndian.Uint64(lenBs)

		reqBs := make([]byte, length)

		_, err = conn.Read(reqBs)
		if err != nil {
			return err
		}

		respData := handleMsg(reqBs)
		respLen := len(respData)

		res := make([]byte, respLen+numOfLengthBytes)

		// 写入长度字段
		binary.BigEndian.PutUint64(res[:numOfLengthBytes], uint64(respLen))
		// 写入响应数据
		copy(res[numOfLengthBytes:], respData)

		_, err = conn.Write(res)
		if err != nil {
			return err
		}
	}
}

func handleMsg(req []byte) []byte {
	res := make([]byte, len(req)*2)
	copy(res[:len(req)], req)
	copy(res[len(req):], req)
	return res
}
