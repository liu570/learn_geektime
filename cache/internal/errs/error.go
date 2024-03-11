package errs

import (
	"errors"
	"fmt"
)

var (
	ErrKeyNotFound = errors.New("cache: 未找到对应值")
)

func NewErrKeyNotFound(key string) error {
	return errors.New(fmt.Sprintf("cache: 未发现 key:%s", key))
}
