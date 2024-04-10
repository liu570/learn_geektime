package message

import (
	"encoding/binary"
)

type Response struct {
	// 协议头长度
	HeaderLength uint32
	// 协议体长度
	BodyLength uint32
	// 消息ID
	MesssageID uint32
	// 版本
	Version uint8
	// 压缩算法
	Compression uint8
	// 序列化协议
	Serialize uint8
	Error     []byte
	Data      []byte
}

func CalRespHeaderLength(resp *Response) {
	headerLength := 15 + len(resp.Error)
	resp.HeaderLength = uint32(headerLength)
}

func CalRespBodyLength(req *Response) {
	req.BodyLength = uint32(len(req.Data))
}

func EncodeResp(resp *Response) []byte {
	CalRespHeaderLength(resp)
	CalRespBodyLength(resp)
	bs := make([]byte, resp.BodyLength+resp.HeaderLength)
	cur := bs
	binary.BigEndian.PutUint32(cur[:4], resp.HeaderLength)
	cur = cur[4:]
	binary.BigEndian.PutUint32(cur[:4], resp.BodyLength)
	cur = cur[4:]
	binary.BigEndian.PutUint32(cur[:4], resp.MesssageID)
	cur = cur[4:]
	cur[0] = resp.Version
	cur[1] = resp.Compression
	cur[2] = resp.Serialize
	cur = cur[3:]
	copy(cur[:len(resp.Error)], resp.Error)
	cur = cur[len(resp.Error):]
	copy(cur, resp.Data)

	return bs
}

func DecodeResp(bs []byte) *Response {
	resp := new(Response)
	resp.HeaderLength = binary.BigEndian.Uint32(bs[:4])
	header := bs[:resp.HeaderLength]
	header = header[4:]
	resp.BodyLength = binary.BigEndian.Uint32(header[:4])
	header = header[4:]
	resp.MesssageID = binary.BigEndian.Uint32(header[:4])
	header = header[4:]
	resp.Version = header[0]
	resp.Compression = header[1]
	resp.Serialize = header[2]
	header = header[3:]
	if len(header) > 0 {
		resp.Error = header
	}
	if len(bs[resp.HeaderLength:]) > 0 {
		resp.Data = bs[resp.HeaderLength:]
	}
	return resp
}
