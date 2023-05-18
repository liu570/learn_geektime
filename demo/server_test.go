package demo

import (
	"encoding/binary"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	//开始监听端口
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		t.Fatal(err)
	}
	// 添加一个 for 循环不然只能接收到一个连接
	for {
		// 开始接收连接
		conn, err := listener.Accept()
		if err != nil {
			t.Fatal(err)
		}

		// 接收到了连接之后交给一个goroutine去处理
		go func() {
			handle(conn)
		}()

	}
}

func handle(conn net.Conn) {
	for {
		lenBs := make([]byte, 8)
		_, err := conn.Read(lenBs)
		if err != nil {
			// 简单处理如果读取数据失败，则直接关闭连接，不进行相应处理
			conn.Close()
			return
		}
		//根据头部八个字节的长度获取需要读取数据的长度再根据获取的长度读取数据
		msgLen := binary.BigEndian.Uint64(lenBs) //解码获取请求的长度
		reqBs := make([]byte, msgLen)
		_, err = conn.Read(reqBs)
		if err != nil {
			conn.Close()
			return
		}
		_, err = conn.Write([]byte("hello, world"))
		if err != nil {
			// 写错误，同理直接关闭连接
			conn.Close()
			return
		}
	}
}
