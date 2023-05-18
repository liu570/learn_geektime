package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"testing"
)

func TestJsonColumn(t *testing.T) {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/test")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.ExecContext(context.Background(),
		`CREATE TABLE IF NOT EXISTS user_tab(
    id INTEGER PRIMARY KEY,
    address TEXT NOT NULL
	)`)
	if err != nil {
		t.Fatal(err)
	}

	res, err := db.Exec("INSERT INTO `user_tab`(`id`,`address`) VALUES(?,?)",
		1, JsonColumn[Address]{Val: Address{Province: "广东", City: "东莞"}})
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

	row := db.QueryRowContext(context.Background(),
		"SELECT `id`,`address` FROM `user_tab` LIMIT ?", 1)
	if row.Err() != nil {
		t.Fatal(row.Err())
	}

	user := &User{}
	err = row.Scan(&user.Id, &user.Address)
	if err != nil {
		t.Fatal(err)
	}
	//user.Address.Val.City
	_, err = db.Exec(`DELETE FROM user_tab;`)
	if err != nil {
		t.Fatal(err)
	}

}

type User struct {
	Id      int64
	Address JsonColumn[Address]
}

type Address struct {
	Province string
	City     string
}

type JsonColumn[T any] struct {
	Val   T
	Valid bool // 标记上数据库存的字段是不是 NULL
}

// 用于查询获取数据库的值，该方法使用指针是因为其要修改结构体的值
func (j *JsonColumn[T]) Scan(src any) error {
	if src == nil {
		return nil
	}
	bs := src.([]byte)
	if len(bs) == 0 {
		return nil
	}
	if err := json.Unmarshal(bs, &j.Val); err != nil {
		return err
	}
	j.Valid = true
	return nil
}

// 用于查询传入数据库的参数
func (j JsonColumn[T]) Value() (driver.Value, error) {
	return json.Marshal(j.Val)
}
