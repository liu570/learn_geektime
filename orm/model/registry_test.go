package model

import (
	"github.com/stretchr/testify/assert"
	"learn_geektime/orm/internal/errs"
	"reflect"
	"testing"
)

func Test_parseModel(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    *Model
		wantErr error
	}{
		{
			name:  "ptr",
			input: &TestModel{},
			want:  getModelTest(),
		},
		{
			name:    "struct",
			input:   TestModel{},
			wantErr: errs.ErrPointOnly,
		},
		{
			name:    "nil",
			input:   nil,
			wantErr: errs.ErrInputNil,
		},
		{
			name:  "Type with nil",
			input: (*TestModel)(nil),
			want:  getModelTest(),
		},
		{
			name:  "column tag",
			input: &ColumnTag{},
			want: &Model{
				TableName: "column_tag",
				Fields: []*Field{
					{
						GoName:  "ID",
						ColName: "id3",
						Type:    reflect.TypeOf(ColumnTag{}).Field(0).Type,
						Offset:  reflect.TypeOf(ColumnTag{}).Field(0).Offset,
						Index:   0,
					},
				},
				FieldMap: map[string]*Field{
					"ID": {
						GoName:  "ID",
						ColName: "id3",
						Type:    reflect.TypeOf(ColumnTag{}).Field(0).Type,
						Offset:  reflect.TypeOf(ColumnTag{}).Field(0).Offset,
						Index:   0,
					},
				},
				ColumnMap: map[string]*Field{
					"id3": {
						GoName:  "ID",
						ColName: "id3",
						Type:    reflect.TypeOf(ColumnTag{}).Field(0).Type,
						Offset:  reflect.TypeOf(ColumnTag{}).Field(0).Offset,
						Index:   0,
					},
				},
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			r := NewRegistry()
			m, err := r.Register(tt.input)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, m)
		})

	}
}

func getModelTest() *Model {
	t := &TestModel{}
	typ := reflect.TypeOf(t).Elem()
	return &Model{
		TableName: "test_model",
		Fields: []*Field{
			{
				GoName:  "Id",
				ColName: "id",
				Type:    typ.Field(0).Type,
				Offset:  typ.Field(0).Offset,
				Index:   0,
			},
			{
				GoName:  "FirstName",
				ColName: "first_name",
				Type:    typ.Field(1).Type,
				Offset:  typ.Field(1).Offset,
				Index:   1,
			},
			{
				GoName:  "Age",
				ColName: "age",
				Type:    typ.Field(2).Type,
				Offset:  typ.Field(2).Offset,
				Index:   2,
			},
			{
				GoName:  "LastName",
				ColName: "last_name",
				Type:    typ.Field(3).Type,
				Offset:  typ.Field(3).Offset,
				Index:   3,
			},
		},
		FieldMap: map[string]*Field{
			"Id": {
				GoName:  "Id",
				ColName: "id",
				Type:    typ.Field(0).Type,
				Offset:  typ.Field(0).Offset,
				Index:   0,
			},
			"FirstName": {
				GoName:  "FirstName",
				ColName: "first_name",
				Type:    typ.Field(1).Type,
				Offset:  typ.Field(1).Offset,
				Index:   1,
			},
			"Age": {
				GoName:  "Age",
				ColName: "age",
				Type:    typ.Field(2).Type,
				Offset:  typ.Field(2).Offset,
				Index:   2,
			},
			"LastName": {
				GoName:  "LastName",
				ColName: "last_name",
				Type:    typ.Field(3).Type,
				Offset:  typ.Field(3).Offset,
				Index:   3,
			},
		},
		ColumnMap: map[string]*Field{
			"id": {
				GoName:  "Id",
				ColName: "id",
				Type:    typ.Field(0).Type,
				Offset:  typ.Field(0).Offset,
				Index:   0,
			},
			"first_name": {
				GoName:  "FirstName",
				ColName: "first_name",
				Type:    typ.Field(1).Type,
				Offset:  typ.Field(1).Offset,
				Index:   1,
			},
			"age": {
				GoName:  "Age",
				ColName: "age",
				Type:    typ.Field(2).Type,
				Offset:  typ.Field(2).Offset,
				Index:   2,
			},
			"last_name": {
				GoName:  "LastName",
				ColName: "last_name",
				Type:    typ.Field(3).Type,
				Offset:  typ.Field(3).Offset,
				Index:   3,
			},
		},
	}
}

type ColumnTag struct {
	ID uint64 `orm:"column=id3"`
}

func NewTesModel() *TestModel {
	return &TestModel{}
}
