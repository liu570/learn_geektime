package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeDecodeReq(t *testing.T) {
	testCases := []struct {
		name string
		req  *Request
	}{
		{
			name: "empty case",
			req:  &Request{},
		},
		{
			name: "no meta and no servername",
			req: &Request{
				MesssageID:  3,
				Version:     1,
				Compression: 1,
				Serializer:  3,
			},
		},
		{
			name: "no meta",
			req: &Request{
				MesssageID:  3,
				Version:     1,
				Compression: 1,
				Serializer:  3,
				ServiceName: "user-service",
				MethodName:  "GetById",
			},
		},
		{
			name: "multi meta",
			req: &Request{
				MesssageID:  3,
				Version:     1,
				Compression: 1,
				Serializer:  3,
				ServiceName: "user-service",
				MethodName:  "GetById",
				Meta: map[string]string{
					"Id":   "123",
					"Addr": "8081",
				},
			},
		},
		{
			name: "meta with no servername",
			req: &Request{
				MesssageID:  3,
				Version:     1,
				Compression: 1,
				Serializer:  3,

				MethodName: "GetById",
				Meta: map[string]string{
					"Id":   "123",
					"Addr": "8081",
				},
			},
		},
		{
			name: "meta with no MethodName",
			req: &Request{
				MesssageID:  3,
				Version:     1,
				Compression: 1,
				Serializer:  3,
				ServiceName: "user-service",
				Meta: map[string]string{
					"Id":   "123",
					"Addr": "8081",
				},
			},
		},
		{
			name: "all in",
			req: &Request{
				MesssageID:  3,
				Version:     1,
				Compression: 1,
				Serializer:  3,
				ServiceName: "user-service",
				MethodName:  "GetById",
				Meta: map[string]string{
					"Id":   "123",
					"Addr": "8081",
				},
				Data: []byte("123456789"),
			},
		},
		{
			name: "data with \\n",
			req: &Request{
				MesssageID:  3,
				Version:     1,
				Compression: 1,
				Serializer:  3,
				ServiceName: "user-service",
				MethodName:  "GetById",
				Meta: map[string]string{
					"Id":   "123",
					"Addr": "8081",
				},
				Data: []byte("123\n456789"),
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			bs := EncodeReq(tt.req)
			req := DecodeReq(bs)
			assert.Equal(t, tt.req, req)
		})
	}
}
