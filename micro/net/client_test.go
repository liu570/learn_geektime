package net

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClient_Send(t *testing.T) {
	server := &Server{
		network: "tcp",
		addr:    ":8081",
	}
	go func() {
		if err := server.Start(); err != nil {
			t.Log(err)
		}
	}()
	time.Sleep(time.Second * 1)

	client := &Client{
		network: "tcp",
		addr:    "localhost:8081",
	}

	testCases := []struct {
		name string

		client   Client
		sendData string

		wantErr  error
		wantResp string
	}{
		{
			name:     "hello",
			sendData: "hello",
			wantResp: "hellohello",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.Send(tc.sendData)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}
