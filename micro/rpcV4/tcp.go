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
	header := binary.BigEndian.Uint32(lenBs[:4])
	body := binary.BigEndian.Uint32(lenBs[4:])
	data := make([]byte, body+header)
	_, err = conn.Read(data[8:])
	copy(data[:8], lenBs)
	return data, err
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
