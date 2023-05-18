package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type sqlTestSuite struct {
	// 测试套件 继承suite.suite
	suite.Suite

	// 配置字段
	// 数据库驱动
	driver string
	dsn    string

	// 初始化字段
	db *sql.DB
}

// 测试套件初始化函数
func (s *sqlTestSuite) SetupSuite() {
	db, err := sql.Open(s.driver, s.dsn)
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db

	// 设置超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// 创建一个名为test_model的表
	_, err = s.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
	)`)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *sqlTestSuite) TearDownSuite() {
	_, err := s.db.Exec(`DELETE FROM test_model;`)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *sqlTestSuite) TestCRUD() {
	// 打开连接
	t := s.T()
	db, err := sql.Open(s.driver, s.dsn)
	if err != nil {
		t.Fatal(err)
	}
	//设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// 执行增删改
	res, err := db.ExecContext(ctx,
		"INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(1,'Tom',18,'Jerry') ")
	if err != nil {
		t.Fatal(err)
	}

	//我们在执行 INSERT、UPDATE 或 DELETE 等写操作之后都需要检查受影响的行数，以确保操作的正确性和一致性。
	//由于每种数据库的实现方式不同，因此需要根据具体情况使用相应的库和方法进行查询和结果处理。
	afferted, err := res.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if afferted != 1 {
		t.Fatal(err)
	}

	// QueryContext 执行有返回行的语句 通常是查询语句 SELECT
	rows, err := db.QueryContext(context.Background(),
		"SELECT `id`,`first_name`,`age`,`last_name` FROM `test_model` LIMIT ?", 1)
	if err != nil {
		t.Fatal(err)
	}
	// 将获取到的行 逐行的写入到TestModel类中
	for rows.Next() {
		tm := &TestModel{}
		err = rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "Tom", tm.FirstName)
	}
	rows.Close()

	//或者执行 Exec(xxx)
	res, err = db.ExecContext(ctx,
		"UPDATE `test_model` SET `first_name` = 'changed' where `id` = ?", 1)
	if err != nil {
		t.Fatal(err)
	}

	afferted, err = res.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if afferted != 1 {
		t.Fatal(err)
	}

	row := db.QueryRowContext(context.Background(),
		"SELECT `id`,`first_name`,`age`,`last_name` FROM `test_model` LIMIT ?", 1)
	tm := &TestModel{}
	err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "changed", tm.FirstName)
}

func (s *sqlTestSuite) TestJsonColumn() {
	t := s.T()
	db, err := sql.Open(s.driver, s.dsn)
	if err != nil {
		t.Fatal(err)
	}
	res, err := db.Exec("INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`)"+
		" VALUES(?,?,?,?)", 1, FullName{FirstName: "A", LastName: "B"}, 18, "Jerry")
	if err != nil {
		t.Fatal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if affected != 1 {
		t.Fatal(err)
	}

}

// 如果要在sql查询的时候使用复杂的参数的话，则需要使该复杂参数（结构体）实现value接口
type FullName struct {
	FirstName string
	LastName  string
}

func (f *FullName) Value() (driver.Value, error) {
	return f.FirstName + f.LastName, nil
}

func TestSQLite(t *testing.T) {
	suite.Run(t, &sqlTestSuite{
		driver: "mysql",
		dsn:    "root:123456@tcp(localhost:3306)/test",
		// 这里root 是数据库用户名、123456 是数据库密码、 test 是对应的数据库名字
	})
}
func TestTimer(t *testing.T) {
	timer := time.NewTimer(0)
	fmt.Println(timer.Stop())
	<-timer.C
}

type TestModel struct {
	Id        int64 `eorm:"auto_increment,primary_key"`
	FirstName string
	Age       string
	LastName  *sql.NullString
}
