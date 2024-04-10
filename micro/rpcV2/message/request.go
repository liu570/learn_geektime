package message

import (
	"bytes"
	"encoding/binary"
)

// Request rpc 调用所需要的调用信息
// warning: 头部字段不允许出现 `\n` 与 `\r` 特殊字段
type Request struct {
	// 协议头长度
	HeaderLength uint32
	// 协议体长度
	BodyLength uint32
	// 消息 ID
	MesssageID uint32
	// 版本
	Version uint8
	// 压缩算法
	Compression uint8
	// 序列化协议
	Serialize uint8
	// 服务名与方法名
	ServiceName string
	MethodName  string
	// 元数据
	Meta map[string]string
	//协议体
	Data []byte
}

func CalReqHeaderLength(req *Request) {
	headerLength := 15 + len(req.MethodName) + 1 + len(req.ServiceName) + 1
	if len(req.Meta) > 0 {
		for key, value := range req.Meta {
			headerLength += len(key) + len(value) + 1
			headerLength += 1
		}
	}
	req.HeaderLength = uint32(headerLength)
}

func CalReqBodyLength(req *Request) {
	req.BodyLength = uint32(len(req.Data))
}

func EncodeReq(req *Request) []byte {
	CalReqHeaderLength(req)
	CalReqBodyLength(req)
	bs := make([]byte, req.BodyLength+req.HeaderLength)
	cur := bs
	binary.BigEndian.PutUint32(cur[:4], req.HeaderLength)
	cur = cur[4:]
	binary.BigEndian.PutUint32(cur[:4], req.BodyLength)
	cur = cur[4:]
	binary.BigEndian.PutUint32(cur[:4], req.MesssageID)
	cur = cur[4:]
	cur[0] = req.Version
	cur[1] = req.Compression
	cur[2] = req.Serialize
	cur = cur[3:]
	copy(cur[0:len(req.ServiceName)], req.ServiceName)
	cur = cur[len(req.ServiceName):]
	cur[0] = '\n'
	cur = cur[1:]
	copy(cur[0:len(req.MethodName)], req.MethodName)
	cur = cur[len(req.MethodName):]
	cur[0] = '\n'
	cur = cur[1:]
	if len(req.Meta) > 0 {
		for key, value := range req.Meta {
			copy(cur[0:len(key)], key)
			cur = cur[len(key):]
			cur[0] = '\r'
			cur = cur[1:]
			copy(cur[0:len(value)], value)
			cur = cur[len(value):]
			cur[0] = '\n'
			cur = cur[1:]
		}
	}
	copy(cur[0:], req.Data)

	return bs
}

func DecodeReq(bs []byte) *Request {
	req := new(Request)
	req.HeaderLength = binary.BigEndian.Uint32(bs[:4])
	header := bs[:req.HeaderLength]
	header = header[4:]
	req.BodyLength = binary.BigEndian.Uint32(header[:4])
	header = header[4:]
	req.MesssageID = binary.BigEndian.Uint32(header[:4])
	header = header[4:]
	req.Version = header[0]
	req.Compression = header[1]
	req.Serialize = header[2]
	header = header[3:]
	index := bytes.IndexByte(header, '\n')
	req.ServiceName = string(header[:index])
	header = header[index+1:]
	index = bytes.IndexByte(header, '\n')
	req.MethodName = string(header[:index])
	header = header[index+1:]
	if len(header) > 0 {
		meta := make(map[string]string, 4)
		for len(header) > 0 {
			index = bytes.IndexByte(header, '\n')
			idx := bytes.IndexByte(header, '\r')
			meta[string(header[:idx])] = string(header[idx+1 : index])
			header = header[index+1:]
		}
		req.Meta = meta
	}
	if req.BodyLength > 0 {
		req.Data = bs[req.HeaderLength:]
	}

	return req
}
