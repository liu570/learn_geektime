package rpc

import (
	"encoding/binary"
	"net"
)

func ReadMsg(conn net.Conn) ([]byte, error) {
	// 响应有多长
	lenBs := make([]byte, numOfLengthBytes)
	_, err := conn.Read(lenBs)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint64(lenBs)
	respBs := make([]byte, length)

	_, err = conn.Read(respBs)
	if err != nil {
		return nil, err
	}
	return respBs, nil
}

func EncodeMsg(data []byte) []byte {
	respLen := len(data)
	res := make([]byte, respLen+numOfLengthBytes)
	// 写入长度字段
	binary.BigEndian.PutUint64(res[:numOfLengthBytes], uint64(respLen))
	// 写入响应数据
	copy(res[numOfLengthBytes:], data)
	return res
}
