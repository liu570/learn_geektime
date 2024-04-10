package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeDecodeResp(t *testing.T) {
	testCases := []struct {
		name string
		resp *Response
	}{
		{
			name: "empty",
			resp: &Response{},
		},
		{
			name: "no data",
			resp: &Response{
				MesssageID:  23,
				Version:     1,
				Compression: 3,
				Serialize:   12,
				Data:        []byte("hello world"),
			},
		},
		{
			name: "error with data",
			resp: &Response{
				MesssageID:  23,
				Version:     1,
				Compression: 3,
				Serialize:   12,
				Error:       []byte("error!"),
				Data:        []byte("hello world"),
			},
		},
		{
			name: "success",
			resp: &Response{
				MesssageID:  23,
				Version:     1,
				Compression: 3,
				Serialize:   12,
				Data:        []byte("hello world"),
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			bs := EncodeResp(tt.resp)
			resp := DecodeResp(bs)
			assert.Equal(t, tt.resp, resp)
		})
	}
}
