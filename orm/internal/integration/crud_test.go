// 下列代表 我们是端到端的测试,是需要依赖外部环境的

//go:build e2e

package integration

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"learn_geektime/orm"
	"testing"
)

type TestSuite struct {
	suite.Suite
	db *orm.DB

	driver string
	dsn    string
}

// 类似于init函数
func (i *TestSuite) SetupSuite() {
	db, err := orm.Open(i.driver, i.dsn)
	if err != nil {
		i.T().Fatal(err)
	}
	i.db = db
	i.db.Wait()
}
func (i *TestSuite) TestInsert() {
	t := i.T()
	db := i.db
	testCases := []struct {
		name string
		i    *orm.Inserter[TestModel]
		// 受影响行数
		affected int64
		wantErr  error
		wantData *TestModel
	}{
		{
			name: "insert single",
			i: orm.NewInserter[TestModel](db).Values(
				&TestModel{
					Id:       13,
					LastName: &sql.NullString{String: "liu", Valid: true},
				},
			),
			affected: 1,
			wantData: &TestModel{
				Id:       13,
				LastName: &sql.NullString{String: "liu", Valid: true},
			},
		},
		{
			name: "insert single",
			i: orm.NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        14,
					FirstName: "WenJie",
					Age:       22,
					LastName:  &sql.NullString{String: "Liu", Valid: true},
				},
			),
			affected: 1,
			wantData: &TestModel{
				Id:        14,
				FirstName: "WenJie",
				Age:       22,
				LastName:  &sql.NullString{String: "Liu", Valid: true},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.i.Exec(context.Background())
			affected, err := res.RowsAffected()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.affected, affected)
			id, err := res.LastInsertId()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			data, err := orm.NewSelector[TestModel](db).Where(orm.C("Id").EQ(id)).Get(context.Background())

			require.NoError(t, err)
			assert.Equal(t, tc.wantData, data)
		})
	}
}

func (u *TestSuite) TestUpdate() {
	t := u.T()
	db := u.db
	testCases := []struct {
		name string
		u    *orm.Updater[TestModel]
		// 受影响行数
		affected int64
		wantErr  error
		wantData *TestModel
	}{
		{
			name: "update single",
			u: orm.NewUpdater[TestModel](db).Update(&TestModel{
				LastName: &sql.NullString{String: "update", Valid: true},
			}).Set(orm.C("LastName")).Where(orm.C("Id").EQ(13)),
			affected: 1,
			wantData: &TestModel{
				Id:       13,
				LastName: &sql.NullString{String: "update", Valid: true},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.u.Exec(context.Background())
			affected, err := res.RowsAffected()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.affected, affected)
			data, err := orm.NewSelector[TestModel](db).Where(orm.C("Id").EQ(13)).Get(context.Background())
			require.NoError(t, err)
			assert.Equal(t, tc.wantData, data)
		})
	}
}

func TestMySQL(t *testing.T) {

	suite.Run(t, &TestSuite{
		driver: "mysql",
		dsn:    "root:root@tcp(127.0.0.1:13306)/integration_test",
	})
}

//func TestSQLite(t *testing.T) {
//	suite.Run(t, &TestSuite{
//		driver: "sqlite3",
//		dsn:    "file:test.db?cache=shared&mode=memory",
//	})
//}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
