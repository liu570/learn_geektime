package orm

import "database/sql"

// Result 的封装 处理 sql.result 的错误 防止我们在平常使用时要处理两次 err
type Result struct {
	res sql.Result
	err error
}

func (r Result) LastInsertId() (int64, error) {
	if r.err != nil {
		return 0, r.err
	}
	return r.res.LastInsertId()
}

func (r Result) RowsAffected() (int64, error) {
	if r.err != nil {
		return 0, r.err
	}
	return r.res.RowsAffected()
}

func (r Result) Err() error {
	return r.err
}
