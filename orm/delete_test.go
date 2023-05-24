package orm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// DELETE Multi规则
// 1.如果你为一个表声明了别名，当你指向这个表的时候，就必须使用这个别名，例如：
//		-- 正确的写法:
//		DELETE t1 FROM test AS t1, test2 WHERE ...
//		-- 错误的写法:
//		DELETE test FROM test AS t1, test2 WHERE ...
// 2.在多个表联合删除时，不能使用 order by 或 limit，而单个表的删除时就没有这个限制。
// 3.当前我们还不能在删除表的时候，在子查询中select from相同的表。

func TestDeleter_Build(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name    string
		d       QueryBuilder
		want    *Query
		wantErr error
	}{
		{
			name: "single delete",
			d:    NewDeleter[TestModel](db),
			want: &Query{
				SQL: "DELETE FROM `test_model`",
			},
			wantErr: nil,
		},
		{
			name: "single delete where",
			d:    NewDeleter[TestModel](db).Where(C("Id").EQ(14)),
			want: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `id` = ?",
				Args: []any{14},
			},
			wantErr: nil,
		},
		{
			name: "single delete order by",
			d:    NewDeleter[TestModel](db).OrderBy(C("Id")),
			want: &Query{
				SQL: "DELETE FROM `test_model` ORDER BY `id`",
			},
			wantErr: nil,
		},
		{
			name: "single delete total",
			d:    NewDeleter[TestModel](db).Where(C("Id").EQ(13)).OrderBy(C("LastName")).Limit(1),
			want: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `id` = ? ORDER BY `last_name` LIMIT ?",
				Args: []any{13, 1},
			},
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.d.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.want, query)
		})
	}
}

//func TestDeleter_MultipleTable(t *testing.T) {
//	db := memoryDB(t)
//
//	type Order struct {
//		Id        int
//		UsingCol1 string
//		UsingCol2 string
//	}
//
//	type OrderDetail struct {
//		OrderId int
//		ItemId  int
//
//		UsingCol1 string
//		UsingCol2 string
//	}
//	type Item struct {
//		Id int
//	}
//	testCases := []struct {
//		name    string
//		d       QueryBuilder
//		want    *Query
//		wantErr error
//	}{
//
//		{
//			name: "multiple delete",
//			d: func() QueryBuilder {
//				t1 := TableOf(&Order{}).As("t1")
//				t2 := TableOf(&OrderDetail{})
//				return NewDeleter[Order](db).
//					From(t1.Join(t2).On(t1.C("Id").EQ(t2.C("OrderId")))).Where(t1.C("Id").EQ(13))
//			}(),
//			want: &Query{
//				SQL:  "DELETE `t1` FROM `test_model` AS `t1`, `test2` WHERE `t1`.`id` = ?",
//				Args: []any{13},
//			},
//			wantErr: nil,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			query, err := tc.d.Build()
//			assert.Equal(t, tc.wantErr, err)
//			if err != nil {
//				return
//			}
//			assert.Equal(t, tc.want, query)
//		})
//	}
//}
