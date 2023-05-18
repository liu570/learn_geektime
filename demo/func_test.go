package demo

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIterateFuncs(t *testing.T) {
	type args struct {
		val any
	}
	tests := []struct {
		//测试名字
		name string
		//测试属性
		args args
		//想要的返回值
		want    map[string]*FuncInfo
		wantErr error
	}{
		{
			name:    "nil",
			wantErr: errors.New("输入 nil"),
		},
		{
			name: "basic type",
			args: args{
				val: 123,
			},
			wantErr: errors.New("不支持类型"),
		},
		{
			name: "struct type",
			args: args{
				val: Order{
					buyer:  18,
					seller: 100,
				},
			},
			want: map[string]*FuncInfo{
				"GetBuyer": {
					Name:   "GetBuyer",
					In:     []reflect.Type{reflect.TypeOf(Order{})},
					Out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(18)},
				},
			},
		},
		{
			name: "struct type but input ptr",
			args: args{
				val: &Order{
					buyer:  18,
					seller: 100,
				},
			},
			want: map[string]*FuncInfo{
				"GetBuyer": {
					Name:   "GetBuyer",
					In:     []reflect.Type{reflect.TypeOf(&Order{})},
					Out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(18)},
				},
			},
		},
		{
			name: "pointer type",
			args: args{
				val: &OrderV1{
					buyer:  18,
					seller: 100,
				},
			},
			want: map[string]*FuncInfo{
				"GetBuyer": {
					Name:   "GetBuyer",
					In:     []reflect.Type{reflect.TypeOf(&OrderV1{})},
					Out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(18)},
				},
			},
		},
		{
			name: "pointer type but input struct",
			args: args{
				val: OrderV1{
					buyer:  18,
					seller: 100,
				},
			},
			want: map[string]*FuncInfo{
				"GetBuyer": {
					Name:   "GetBuyer",
					In:     []reflect.Type{reflect.TypeOf(&OrderV1{})},
					Out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(18)},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IterateFuncs(tt.args.val)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

type Order struct {
	buyer  int64
	seller int64
}

// 反射层面认为函数是 GetBuyer(o Order)
func (o Order) GetBuyer() int64 {
	return o.buyer
}

//func (o Order) getSeller() int64 {
//	return o.seller
//}

type OrderV1 struct {
	buyer  int64
	seller int64
}

// 反射层面认为函数是 GetBuyer(o Order)
func (o *OrderV1) GetBuyer() int64 {
	return o.buyer
}

type MyInterface interface {
	Abc()
}

// 这句用来确认abcImpl 确实实现了 MyInterface 接口
var _ MyInterface = &abcImpl{}

type abcImpl struct {
}

func (a *abcImpl) Abc() {
	//TODO implement me
	panic("implement me")
}
