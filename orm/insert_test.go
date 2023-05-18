package orm

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"learn_geektime/orm/internal/errs"
	"testing"
)

func TestInserter_Build_Unsafe(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name    string
		i       QueryBuilder
		want    *Query
		wantErr error
	}{
		{
			name:    "no value",
			i:       NewInserter[TestModel](db).Values(),
			wantErr: errs.ErrInsertZeroRow,
		},
		{
			// 插入一行
			name: "simple value",
			i: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        12,
					FirstName: "Tom",
					Age:       18,
					LastName:  &sql.NullString{Valid: true, String: "Jerry"},
				}),
			want: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?);",
				Args: []any{
					int64(12), "Tom", int8(18), &sql.NullString{Valid: true, String: "Jerry"},
				},
			},
		},
		{
			// 插入多行
			name: "multi value",
			i: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        12,
					FirstName: "Tom",
					Age:       18,
					LastName:  &sql.NullString{Valid: true, String: "Jerry"},
				},
				&TestModel{
					FirstName: "Deng",
					Id:        13,
					Age:       17,
					LastName:  &sql.NullString{Valid: true, String: "Ming"},
				}),
			want: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?),(?,?,?,?);",
				Args: []any{
					int64(12), "Tom", int8(18), &sql.NullString{Valid: true, String: "Jerry"},
					int64(13), "Deng", int8(17), &sql.NullString{Valid: true, String: "Ming"},
				},
			},
		},
		{
			// 指定列名字
			name: "specify columns",
			i: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        12,
					FirstName: "Tom",
					Age:       18,
					LastName:  &sql.NullString{Valid: true, String: "Jerry"},
				},
				&TestModel{
					FirstName: "Deng",
					Id:        13,
					Age:       17,
					LastName:  &sql.NullString{Valid: true, String: "Ming"},
				}).Columns("Id", "Age"),
			want: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`age`) VALUES(?,?),(?,?);",
				Args: []any{
					int64(12), int8(18),
					int64(13), int8(17),
				},
			},
		},
		{
			// update or insert
			name: "upset",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			}).OnConflict().Update(Assign("Age", 19)),
			want: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?)" +
					" ON DUPLICATE KEY UPDATE `age` = ?;",
				Args: []any{
					int64(12), "Tom", int8(18), &sql.NullString{Valid: true, String: "Jerry"}, 19,
				},
			},
		},

		{
			// update or insert
			name: "upset use columns",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			}).OnConflict().Update(C("Age")),
			want: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?)" +
					" ON DUPLICATE KEY UPDATE `age` = VALUES(`age`);",
				Args: []any{
					int64(12), "Tom", int8(18), &sql.NullString{Valid: true, String: "Jerry"},
				},
			},
		},

		{
			// update or insert
			name: "upset multiple",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			}).OnConflict().Update(Assign("Age", 19), C("FirstName")),
			want: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?)" +
					" ON DUPLICATE KEY UPDATE `age` = ?,`first_name` = VALUES(`first_name`);",
				Args: []any{
					int64(12), "Tom", int8(18), &sql.NullString{Valid: true, String: "Jerry"}, 19,
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Build()
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
func TestInserter_Build_Reflect(t *testing.T) {
	db := memoryDBWithReflect(t)
	testCases := []struct {
		name    string
		i       QueryBuilder
		want    *Query
		wantErr error
	}{
		{
			name:    "no value",
			i:       NewInserter[TestModel](db).Values(),
			wantErr: errs.ErrInsertZeroRow,
		},
		{
			// 插入一行
			name: "simple value",
			i: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        12,
					FirstName: "Tom",
					Age:       18,
					LastName:  &sql.NullString{Valid: true, String: "Jerry"},
				}),
			want: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?);",
				Args: []any{
					int64(12), "Tom", int8(18), &sql.NullString{Valid: true, String: "Jerry"},
				},
			},
		},
		{
			// 插入多行
			name: "multi value",
			i: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        12,
					FirstName: "Tom",
					Age:       18,
					LastName:  &sql.NullString{Valid: true, String: "Jerry"},
				},
				&TestModel{
					FirstName: "Deng",
					Id:        13,
					Age:       17,
					LastName:  &sql.NullString{Valid: true, String: "Ming"},
				}),
			want: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?),(?,?,?,?);",
				Args: []any{
					int64(12), "Tom", int8(18), &sql.NullString{Valid: true, String: "Jerry"},
					int64(13), "Deng", int8(17), &sql.NullString{Valid: true, String: "Ming"},
				},
			},
		},
		{
			// 指定列名字
			name: "specify columns",
			i: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        12,
					FirstName: "Tom",
					Age:       18,
					LastName:  &sql.NullString{Valid: true, String: "Jerry"},
				},
				&TestModel{
					FirstName: "Deng",
					Id:        13,
					Age:       17,
					LastName:  &sql.NullString{Valid: true, String: "Ming"},
				}).Columns("Id", "Age"),
			want: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`age`) VALUES(?,?),(?,?);",
				Args: []any{
					int64(12), int8(18),
					int64(13), int8(17),
				},
			},
		},
		{
			// update or insert
			name: "upset",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			}).OnConflict().Update(Assign("Age", 19)),
			want: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?)" +
					" ON DUPLICATE KEY UPDATE `age` = ?;",
				Args: []any{
					int64(12), "Tom", int8(18), &sql.NullString{Valid: true, String: "Jerry"}, 19,
				},
			},
		},

		{
			// update or insert
			name: "upset use columns",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			}).OnConflict().Update(C("Age")),
			want: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?)" +
					" ON DUPLICATE KEY UPDATE `age` = VALUES(`age`);",
				Args: []any{
					int64(12), "Tom", int8(18), &sql.NullString{Valid: true, String: "Jerry"},
				},
			},
		},

		{
			// update or insert
			name: "upset multiple",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			}).OnConflict().Update(Assign("Age", 19), C("FirstName")),
			want: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?)" +
					" ON DUPLICATE KEY UPDATE `age` = ?,`first_name` = VALUES(`first_name`);",
				Args: []any{
					int64(12), "Tom", int8(18), &sql.NullString{Valid: true, String: "Jerry"}, 19,
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Build()
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
func TestInserter_Exec(t *testing.T) {
	//var i *Inserter[TestModel]
	//
	//res := i.Exec(context.Background())
	//id, err := res.LastInsertId()
	//if err != nil {
	//	t.Fatal(err)
	//}

	// 接下来你才能使用这个 id
}
