package valuer

import "testing"

func Test_UnsafeValue_SetColumn(t *testing.T) {
	testSetColumn(t, NewUnsafeValue)
}
