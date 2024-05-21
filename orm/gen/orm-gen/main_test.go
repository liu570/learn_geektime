package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_gen(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	err := gen(buffer, "test_data/user.go")
	require.NoError(t, err)
	assert.Equal(t, `package test_data
import(
    "learn_geektime/orm"
	"database/sql"
)
`, buffer.String())
}
func Test_genFile(t *testing.T) {
	f, err := os.Create("test_data/user.gen.go")
	require.NoError(t, err)
	err = gen(f, "test_data/user.go")
	require.NoError(t, err)
	err = f.Close()
}
