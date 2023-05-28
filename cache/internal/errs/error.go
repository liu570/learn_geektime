package errs

import (
	"errors"
	"fmt"
)

func NewErrKeyNotFound(key string) error {
	return errors.New(fmt.Sprintf("cache: 未发现 key:%s", key))
}
