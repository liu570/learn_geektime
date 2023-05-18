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
				FieldMap: map[string]*Field{
					"ID": {
						GoName:  "ID",
						ColName: "id3",
						Type:    reflect.TypeOf(ColumnTag{}).Field(0).Type,
					},
				},
				ColumnMap: map[string]*Field{
					"id3": {
						GoName:  "ID",
						ColName: "id3",
						Type:    reflect.TypeOf(ColumnTag{}).Field(0).Type,
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
	return &Model{
		TableName: "test_model",
		FieldMap: map[string]*Field{
			"Id": {
				GoName:  "Id",
				ColName: "id",
				Type:    reflect.TypeOf(TestModel{}).Field(0).Type,
			},
			"FirstName": {
				GoName:  "FirstName",
				ColName: "first_name",
				Type:    reflect.TypeOf(TestModel{}).Field(1).Type,
			},
			"Age": {
				GoName:  "Age",
				ColName: "age",
				Type:    reflect.TypeOf(TestModel{}).Field(2).Type,
			},
			"LastName": {
				GoName:  "LastName",
				ColName: "last_name",
				Type:    reflect.TypeOf(TestModel{}).Field(3).Type,
			},
		},
		ColumnMap: map[string]*Field{
			"id": {
				GoName:  "Id",
				ColName: "id",
				Type:    reflect.TypeOf(TestModel{}).Field(0).Type,
			},
			"first_name": {
				GoName:  "FirstName",
				ColName: "first_name",
				Type:    reflect.TypeOf(TestModel{}).Field(1).Type,
			},
			"age": {
				GoName:  "Age",
				ColName: "age",
				Type:    reflect.TypeOf(TestModel{}).Field(2).Type,
			},
			"last_name": {
				GoName:  "LastName",
				ColName: "last_name",
				Type:    reflect.TypeOf(TestModel{}).Field(3).Type,
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
