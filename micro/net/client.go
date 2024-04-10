package net

import (
	"encoding/binary"
	"net"
	"time"
)

type Client struct {
	network string
	addr    string
}

func (c *Client) Send(data string) (string, error) {
	conn, err := net.DialTimeout(c.network, c.addr, time.Second*3)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = conn.Close()
	}()

	reqLen := len(data)

	req := make([]byte, reqLen+numOfLengthBytes)

	// 写入长度字段
	binary.BigEndian.PutUint64(req[:numOfLengthBytes], uint64(reqLen))
	// 写入请求数据
	copy(req[numOfLengthBytes:], data)

	// 发送请求数据
	_, err = conn.Write(req)
	if err != nil {
		return "", err
	}

	// 响应有多长
	lenBs := make([]byte, numOfLengthBytes)
	_, err = conn.Read(lenBs)
	if err != nil {
		return "", err
	}
	length := binary.BigEndian.Uint64(lenBs)
	respBs := make([]byte, length)

	_, err = conn.Read(respBs)
	if err != nil {
		return "", err
	}
	return string(respBs), nil
}
