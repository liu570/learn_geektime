package orm

import (
	"context"
	"database/sql"
	"go.uber.org/multierr"
	"learn_geektime/orm/internal/valuer"
	"learn_geektime/orm/model"
)

type DBOption func(db *DB)

// DB 是sql.DB 的装饰器
type DB struct {
	// core 中包含model.Registry, dialect, valuer.Creator
	core
	db *sql.DB
	//useUnsafe bool //不推荐使用标记位 ， 不好维护
}

func (db *DB) getCore() core {
	return db.core
}
func (db *DB) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.db.QueryContext(ctx, query, args...)
}
func (db *DB) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.db.ExecContext(ctx, query, args...)
}

// Wait 确保我们已经连接上数据库。
// 此时这个方法还并不完善，如果未连接上可能会一直阻塞在这里(所以只能用于测试)
func (db *DB) Wait() {
	for db.db.PingContext(context.Background()) != nil {

	}
}

// Begin 开始一个事务
func (db *DB) Begin(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{
		tx: tx,
	}, nil
}

// DoTx 帮用户管理事务，如果未发生错误，则提交，发生错误则回滚
// 该方法为何要放在 DB 中暂时还不了解
func (db *DB) DoTx(ctx context.Context, opts *sql.TxOptions,
	task func(ctx context.Context, tx *Tx) error) (err error) {
	tx, err := db.Begin(ctx, opts)
	if err != nil {
		return err
	}
	panicked := true
	defer func() {
		if panicked || err != nil {
			err = multierr.Combine(err, tx.Rollback())
		} else {
			err = multierr.Combine(err, tx.Commit())
		}
	}()
	err = task(ctx, tx)
	panicked = false
	return
}

// Open 打开一个 sql.DB 并封装进我们的 DB 结构体
func Open(driver string, dsn string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

// OpenDB
// 我可以用 OpenDB 来传入一个 mock 的 DB
// sqlmock.Open 的 DB
func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		core: core{
			r:          model.NewRegistry(),
			valCreator: valuer.NewUnsafeValue,
			dialect:    DialectMySQL,
		},
		db: db,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

// DBWithRegistry 自己定义 db 里面的 registry
func DBWithRegistry(r model.Registry) DBOption {
	return func(db *DB) {
		db.r = r
	}
}

// DBWithMiddlewares 在 DB 层面上确认 middleware 中间件
func DBWithMiddlewares(ms ...Middleware) DBOption {
	return func(db *DB) {
		db.ms = ms
	}
}

// DBUseReflectValuer 将 db 默认的 unsafe 实现改成 反射 实现
func DBUseReflectValuer() DBOption {
	return func(db *DB) {
		db.valCreator = valuer.NewReflectValue
	}
}

// DBWithDialect 使用该方法表明我们要使用哪个数据库的方言抽象
func DBWithDialect(dialect Dialect) DBOption {
	return func(db *DB) {
		db.dialect = dialect
	}
}
