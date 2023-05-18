package demo

import (
	"encoding/binary"
	"net"
	"testing"
)

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", ":8081")
	if err != nil {
		t.Fatal(err)
	}
	// 写请求
	msg := "how are you"
	msgLen := len(msg)
	// msgLen how are you
	// 数据在电脑中有不同的编码方式，大端编码或者是小端编码
	msgLenBs := make([]byte, 8)
	binary.BigEndian.PutUint64(msgLenBs, uint64(msgLen))
	data := append(msgLenBs, []byte(msg)...) //切片连接切片要在后面加...
	_, err = conn.Write(data)
	if err != nil {
		conn.Close()
		return
	}
	// 读响应
	respBs := make([]byte, 16)
	_, err = conn.Read(respBs)
	if err != nil {
		conn.Close()
	}
}
