package valuer

import (
	"database/sql"
	"learn_geektime/orm/internal/errs"
	"learn_geektime/orm/model"
	"reflect"
)

type reflectValue struct {
	val  reflect.Value
	meta *model.Model
}

var _ Creator = NewReflectValue

func NewReflectValue(val any, meta *model.Model) Value {
	return reflectValue{
		val:  reflect.ValueOf(val).Elem(),
		meta: meta,
	}
}
func (r reflectValue) Field(name string) (any, error) {
	typ := r.val.Type()
	_, ok := typ.FieldByName(name)
	if !ok {
		return nil, errs.NewErrUnknownField(name)
	}
	return r.val.FieldByName(name).Interface(), nil
}

func (r reflectValue) SetColumn(rows *sql.Rows) error {
	if !rows.Next() {
		return errs.ErrNoRows
	}
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	// 如果列过多
	if len(cols) > len(r.meta.FieldMap) {
		return errs.ErrTooManyReturnedColumns
	}
	// 把 cols 映射过去
	colVals := make([]any, 0, len(cols))
	colEleVals := make([]reflect.Value, 0, len(cols))
	for _, col := range cols {
		fd, ok := r.meta.ColumnMap[col]
		if !ok {
			return errs.NewErrUnknownColumn(col)
		}
		fdVal := reflect.New(fd.Type)
		colEleVals = append(colEleVals, fdVal.Elem())
		// Scan 函数传入的值需要指针，所以我们这里不需要调用 Elem
		colVals = append(colVals, fdVal.Interface())
	}

	err = rows.Scan(colVals...)
	if err != nil {
		return err
	}

	for i, col := range cols {
		fd := r.meta.ColumnMap[col]
		//r.val.FieldByName(fd.GoName).Set(colEleVals[i])
		r.val.Field(fd.Index).Set(colEleVals[i])
	}
	return nil
}
