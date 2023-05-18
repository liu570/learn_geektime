package valuer

import (
	"database/sql"
	"fmt"
	"learn_geektime/orm/internal/errs"
	"learn_geektime/orm/model"
	"reflect"
	"unsafe"
)

type unsafeValue struct {
	meta *model.Model
	addr unsafe.Pointer
}

func NewUnsafeValue(val any, meta *model.Model) Value {
	// val 此地传入的是指针 ， 所以我们直接获得他的地址
	addr := unsafe.Pointer(reflect.ValueOf(val).Pointer())
	return unsafeValue{
		addr: addr,
		meta: meta,
	}
}

func (u unsafeValue) Field(name string) (any, error) {
	fdMeta, ok := u.meta.FieldMap[name]
	if !ok {
		return nil, fmt.Errorf("invalid field: %s", name)
	}
	ptr := unsafe.Pointer(uintptr(u.addr) + fdMeta.Offset)
	if ptr == nil {
		return nil, fmt.Errorf("invalid address of the field: %s", name)
	}
	val := reflect.NewAt(fdMeta.Type, ptr)
	return val.Elem().Interface(), nil
}

func (u unsafeValue) SetColumn(rows *sql.Rows) error {
	if !rows.Next() {
		return errs.ErrNoRows
	}
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	if len(cols) > len(u.meta.FieldMap) {
		return errs.ErrTooManyReturnedColumns
	}

	colVals := make([]any, 0, len(cols))
	for _, col := range cols {
		fd, ok := u.meta.ColumnMap[col]
		if !ok {
			return errs.NewErrUnknownColumn(col)
		}
		// 需要计算字段的真实地址： 字段起始地址 + 字段偏移量
		fdVal := reflect.NewAt(fd.Type, unsafe.Pointer(uintptr(u.addr)+fd.Offset))
		colVals = append(colVals, fdVal.Interface())
	}
	return rows.Scan(colVals...)
}
