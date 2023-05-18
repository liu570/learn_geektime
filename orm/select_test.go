package orm

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"learn_geektime/orm/internal/errs"
	"testing"
)

func TestSelector_Join(t *testing.T) {
	db := memoryDB(t)

	type Order struct {
		Id        int
		UsingCol1 string
		UsingCol2 string
	}

	type OrderDetail struct {
		OrderId int
		ItemId  int

		UsingCol1 string
		UsingCol2 string
	}

	type Item struct {
		Id int
	}

	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// 虽然泛型是 Order，但是我们传入 OrderDetail
			name: "specify table",
			q:    NewSelector[Order](db).From(TableOf(&OrderDetail{})),
			wantQuery: &Query{
				SQL: "SELECT * FROM `order_detail`;",
			},
		},
		{
			name: "join",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{})
				return NewSelector[Order](db).
					From(t1.Join(t2).On(t1.C("Id").EQ(t2.C("OrderId"))))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM (`order` AS `t1` JOIN `order_detail` ON `t1`.`id` = `order_id`);",
			},
		},
		{
			name: "multiple join",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{}).As("t2")
				t3 := TableOf(&Item{}).As("t3")
				return NewSelector[Order](db).
					From(t1.Join(t2).
						On(t1.C("Id").EQ(t2.C("OrderId"))).
						Join(t3).On(t2.C("ItemId").EQ(t3.C("Id"))))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM ((`order` AS `t1` JOIN `order_detail` AS `t2` ON `t1`.`id` = `t2`.`order_id`) JOIN `item` AS `t3` ON `t2`.`item_id` = `t3`.`id`);",
			},
		},
		{
			name: "left multiple join",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{}).As("t2")
				t3 := TableOf(&Item{}).As("t3")
				return NewSelector[Order](db).
					From(t1.LeftJoin(t2).
						On(t1.C("Id").EQ(t2.C("OrderId"))).
						LeftJoin(t3).On(t2.C("ItemId").EQ(t3.C("Id"))))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM ((`order` AS `t1` LEFT JOIN `order_detail` AS `t2` ON `t1`.`id` = `t2`.`order_id`) LEFT JOIN `item` AS `t3` ON `t2`.`item_id` = `t3`.`id`);",
			},
		},
		{
			name: "right multiple join",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{}).As("t2")
				t3 := TableOf(&Item{}).As("t3")
				return NewSelector[Order](db).
					From(t1.RightJoin(t2).
						On(t1.C("Id").EQ(t2.C("OrderId"))).
						RightJoin(t3).On(t2.C("ItemId").EQ(t3.C("Id"))))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM ((`order` AS `t1` RIGHT JOIN `order_detail` AS `t2` ON `t1`.`id` = `t2`.`order_id`) RIGHT JOIN `item` AS `t3` ON `t2`.`item_id` = `t3`.`id`);",
			},
		},
		{
			name: "join multiple using",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{})
				return NewSelector[Order](db).
					From(t1.Join(t2).Using("UsingCol1", "UsingCol2"))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM (`order` AS `t1` JOIN `order_detail` USING (`using_col1`,`using_col2`));",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

//func TestSelector_Subquery(t *testing.T) {
//	db := memoryDB(t)
//	type Order struct {
//		Id        int
//		UsingCol1 string
//		UsingCol2 string
//	}
//
//	type OrderDetail struct {
//		OrderId int
//		ItemId  int
//	}
//
//	testCases := []struct {
//		name      string
//		q         QueryBuilder
//		wantQuery *Query
//		wantErr   error
//	}{
//		{
//			name: "from",
//			q: func() QueryBuilder {
//				sub := NewSelector[OrderDetail](db).AsSubquery("sub")
//				return NewSelector[Order](db).From(sub)
//			}(),
//			wantQuery: &Query{
//				SQL: "SELECT * FROM (SELECT * FROM `order_detail`) AS `sub`;",
//			},
//		},
//		//{
//		//	name: "in",
//		//	q: func() QueryBuilder {
//		//		sub := NewSelector[OrderDetail](db).Select(C("OrderId")).AsSubquery("sub")
//		//		return NewSelector[Order](db).Where(C("Id").InQuery(sub))
//		//	}(),
//		//	wantQuery: &Query{
//		//		SQL: "SELECT * FROM `order` WHERE `id` IN (SELECT `order_id` FROM `order_detail`);",
//		//	},
//		//},
//		//{
//		//	name: "exist",
//		//	q: func() QueryBuilder {
//		//		sub := NewSelector[OrderDetail](db).Select(C("OrderId")).AsSubquery("sub")
//		//		return NewSelector[Order](db).Where(Exist(sub))
//		//	}(),
//		//	wantQuery: &Query{
//		//		SQL: "SELECT * FROM `order` WHERE  EXIST (SELECT `order_id` FROM `order_detail`);",
//		//	},
//		//},
//		//{
//		//	name: "not exist",
//		//	q: func() QueryBuilder {
//		//		sub := NewSelector[OrderDetail](db).Select(C("OrderId")).AsSubquery("sub")
//		//		return NewSelector[Order](db).Where(Not(Exist(sub)))
//		//	}(),
//		//	wantQuery: &Query{
//		//		SQL: "SELECT * FROM `order` WHERE  NOT ( EXIST (SELECT `order_id` FROM `order_detail`));",
//		//	},
//		//},
//		//{
//		//	name: "all",
//		//	q: func() QueryBuilder {
//		//		sub := NewSelector[OrderDetail](db).Select(C("OrderId")).AsSubquery("sub")
//		//		return NewSelector[Order](db).Where(C("Id").GT(All(sub)))
//		//	}(),
//		//	wantQuery: &Query{
//		//		SQL: "SELECT * FROM `order` WHERE `id` > ALL (SELECT `order_id` FROM `order_detail`);",
//		//	},
//		//},
//		//{
//		//	name: "some and any",
//		//	q: func() QueryBuilder {
//		//		sub := NewSelector[OrderDetail](db).Select(C("OrderId")).AsSubquery("sub")
//		//		return NewSelector[Order](db).Where(C("Id").GT(Some(sub)), C("Id").LT(Any(sub)))
//		//	}(),
//		//	wantQuery: &Query{
//		//		SQL: "SELECT * FROM `order` WHERE (`id` > SOME (SELECT `order_id` FROM `order_detail`)) AND (`id` < ANY (SELECT `order_id` FROM `order_detail`;SELECT `order_id` FROM `order_detail`));",
//		//	},
//		//},
//	}
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			query, err := tc.q.Build()
//			assert.Equal(t, tc.wantErr, err)
//			if err != nil {
//				return
//			}
//			assert.Equal(t, tc.wantQuery, query)
//		})
//	}
//}

func TestSelector_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	testCases := []struct {
		name     string
		query    string
		mockErr  error
		mockRows *sqlmock.Rows
		wantErr  error
		wantVal  *TestModel
	}{
		{
			name:    "single row",
			query:   "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				rows.AddRow([]byte("123"), []byte("Ming"), []byte("18"), []byte("Deng"))
				return rows
			}(),
			wantVal: &TestModel{
				Id:        123,
				FirstName: "Ming",
				Age:       18,
				LastName: &sql.NullString{
					Valid:  true,
					String: "Deng",
				},
			},
		},
		{
			name:    "invalid columns",
			query:   "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "gender"})
				rows.AddRow([]byte("123"), []byte("Ming"), []byte("male"))
				return rows
			}(),
			wantErr: errs.NewErrUnknownColumn("gender"),
		},
		{
			name:    "less columns",
			query:   "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name"})
				rows.AddRow([]byte("123"), []byte("Ming"))
				return rows
			}(),
			wantVal: &TestModel{
				Id:        123,
				FirstName: "Ming",
			},
		},
	}

	for _, tc := range testCases {
		if tc.mockErr != nil {
			mock.ExpectQuery(tc.query).WillReturnError(tc.mockErr)
		} else {
			mock.ExpectQuery(tc.query).WillReturnRows(tc.mockRows)
		}
	}
	db, err := OpenDB(mockDB, DBUseReflectValuer())
	require.NoError(t, err)
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := NewSelector[TestModel](db).Get(context.Background())
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantVal, res)
		})
	}
}

func TestSelector_Select(t *testing.T) {
	db := memoryDB(t)

	tests := []struct {
		name    string
		s       QueryBuilder
		want    *Query
		wantErr error
	}{
		{
			// 没有指定
			name: "all",
			s:    NewSelector[TestModel](db),
			want: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},

		{
			// 指定列
			name: "specify columns",
			s:    NewSelector[TestModel](db).Select(C("Id"), C("Age")),
			want: &Query{
				SQL: "SELECT `id`,`age` FROM `test_model`;",
			},
		},
		{
			// 指定聚合函数
			// AV,COUNT,SUM,MIN,MAX(xxx)
			name: "specify columns",
			s:    NewSelector[TestModel](db).Select(Min("Id"), Avg("Age")),
			want: &Query{
				SQL: "SELECT MIN(`id`),AVG(`age`) FROM `test_model`;",
			},
		},

		{
			// 特殊的聚合函数
			// 提供一个万金油方法，让用户自己去管对错
			name: "specify aggregate",
			s:    NewSelector[TestModel](db).Select(Raw("DISTINCT `first_name`")),
			want: &Query{
				SQL: "SELECT DISTINCT `first_name` FROM `test_model`;",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Build()
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSelector_Having(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// 调用了，但是啥也没传
			name: "none",
			q:    NewSelector[TestModel](db).GroupBy(C("Age")).Having(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` GROUP BY `age`;",
			},
		},
		{
			// 单个条件
			name: "single",
			q: NewSelector[TestModel](db).GroupBy(C("Age")).
				Having(C("FirstName").EQ("Deng")),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` GROUP BY `age` HAVING `first_name` = ?;",
				Args: []any{"Deng"},
			},
		},
		{
			// 多个条件
			name: "multiple",
			q: NewSelector[TestModel](db).GroupBy(C("Age")).
				Having(C("FirstName").EQ("Deng"), C("LastName").EQ("Ming")),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` GROUP BY `age` HAVING (`first_name` = ?) AND (`last_name` = ?);",
				Args: []any{"Deng", "Ming"},
			},
		},
		{
			// 多个条件
			name: "multiple",
			q: NewSelector[TestModel](db).GroupBy(C("Age")).
				Having(Avg("Age").LT(18), C("FirstName").EQ("Deng"), C("LastName").EQ("Ming")),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` GROUP BY `age` HAVING ((AVG(`age`) < ?) AND (`first_name` = ?)) AND (`last_name` = ?);",
				Args: []any{18, "Deng", "Ming"},
			},
		},
		{
			// 聚合函数
			name: "avg",
			q: NewSelector[TestModel](db).GroupBy(C("Age")).
				Having(Avg("Age").EQ(18)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` GROUP BY `age` HAVING AVG(`age`) = ?;",
				Args: []any{18},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

func TestSelector_GroupBy(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// 调用了，但是啥也没传
			name: "none",
			q:    NewSelector[TestModel](db).GroupBy(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			// 单个
			name: "single",
			q:    NewSelector[TestModel](db).GroupBy(C("Age")),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` GROUP BY `age`;",
			},
		},
		{
			// 多个
			name: "multiple",
			q:    NewSelector[TestModel](db).GroupBy(C("Age"), C("FirstName")),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` GROUP BY `age`,`first_name`;",
			},
		},
		{
			// 不存在
			name:    "invalid column",
			q:       NewSelector[TestModel](db).GroupBy(C("Invalid")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

func TestSelector_OrderBy(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "column",
			q:    NewSelector[TestModel](db).OrderBy(Asc("Age")),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` ORDER BY `age` ASC;",
			},
		},
		{
			name: "columns no order",
			q:    NewSelector[TestModel](db).OrderBy(C("Age"), Desc("Id")),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` ORDER BY `age`,`id` DESC;",
			},
		},
		{
			name: "columns",
			q:    NewSelector[TestModel](db).OrderBy(Asc("Age"), Desc("Id")),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` ORDER BY `age` ASC,`id` DESC;",
			},
		},
		{
			name:    "invalid column",
			q:       NewSelector[TestModel](db).OrderBy(Asc("Invalid")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

func TestSelector_OffsetLimit(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "offset only",
			q:    NewSelector[TestModel](db).Offset(10),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` OFFSET ?;",
				Args: []any{10},
			},
		},
		{
			name: "limit only",
			q:    NewSelector[TestModel](db).Limit(10),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` LIMIT ?;",
				Args: []any{10},
			},
		},
		{
			name: "limit offset",
			q:    NewSelector[TestModel](db).Limit(20).Offset(10),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` LIMIT ? OFFSET ?;",
				Args: []any{20, 10},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

//func TestSelector_Tx(t *testing.T) {
//	db := memoryDB(t)
//	tx, err := db.Begin(context.Background(), &sql.TxOptions{})
//	if err != nil {
//		t.Fatal(err)
//	}
//	// 新建方法传入 db 的话就是脱离这个事务
//	NewSelector[TestModel](tx)
//	//
//	//dao := &UserDAO{
//	//	sess: db.db,
//	//}
//	//txDao := &UserDAO{
//	//	sess: tx.tx,
//	//}
//}

type UserDAO struct {
	sess interface {
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	}
}

func memoryDB(t *testing.T) *DB {
	orm, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		t.Fatal(err)
	}
	return orm
}

func memoryDBWithReflect(t *testing.T) *DB {
	orm, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory", DBUseReflectValuer())
	if err != nil {
		t.Fatal(err)
	}
	return orm
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
