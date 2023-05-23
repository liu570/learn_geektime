package orm

import (
	"context"
	"learn_geektime/orm/internal/errs"
	"learn_geektime/orm/model"
)

type Inserter[T any] struct {
	// 组合了多种语法需要使用的数据
	builder
	// 组合了 model.Register, valuer.Creator, dialect
	core
	// 确定 Inserter 使用的 DB
	sess Session
	// 表明要插入多少个数据 存储多个数据对应的结构体
	values []*T
	// 表明需要插入的列名
	columns []string
	// 表示支持复杂的语句 可以参考 dialect.learn_png 图示理解
	onDuplicate *OnConflictKey
}

func (i *Inserter[T]) Exec(ctx context.Context) Result {
	// 这里调用 Get 方法获取元数据，是怕 Inserter 中的元数据为空
	model, err := i.r.Get(new(T))
	if err != nil {
		return Result{
			err: err,
		}
	}
	qr := exec[T](ctx, i.sess.getCore(), i.sess, &QueryContext{
		Type:      "INSERT",
		Builder:   i,
		Model:     model,
		TableName: model.TableName,
	})
	if qr.Err != nil {
		return Result{
			err: qr.Err,
		}
	}
	return qr.Result.(Result)
}

// NewInserter 创建 insert 构建对象 , 输入参数可以为 DB，Tx
func NewInserter[T any](sess Session) *Inserter[T] {
	c := sess.getCore()
	return &Inserter[T]{
		sess:    sess,
		core:    c,
		builder: builder{dialect: c.dialect},
	}
}

func (i *Inserter[T]) Columns(columns ...string) *Inserter[T] {
	i.columns = columns
	return i
}

// Values  向 insert 阐明我们需要插入多少个对象
func (i *Inserter[T]) Values(vals ...*T) *Inserter[T] {
	i.values = vals
	return i
}

// OnConflict 方言，在此实现不同数据库所特有的方言
func (i *Inserter[T]) OnConflict() *OnConflictBuilder[T] {
	return &OnConflictBuilder[T]{
		i: i,
	}
}

func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}
	var err error
	i.model, err = i.r.Get(i.values[0])
	if err != nil {
		return nil, err
	}

	// INSERT INTO
	i.sb.WriteString("INSERT INTO ")
	i.builder.quote(i.model.TableName)

	// 表明 选定要插入进表的数据
	i.sb.WriteByte('(')
	// 未指定就选择全部字段
	fields := i.model.Fields
	// 指定后就选择指定的字段
	if len(i.columns) != 0 {
		fields = make([]*model.Field, 0, len(i.columns))
		for _, column := range i.columns {
			fd, ok := i.model.FieldMap[column]
			if !ok {
				return nil, errs.NewErrUnknownField(column)
			}
			fields = append(fields, fd)
		}
	}

	// 讲选定的字段编辑进 SQL 语句中去
	for k, field := range fields {
		if k > 0 {
			i.sb.WriteByte(',')
		}
		i.builder.quote(field.ColName)
	}
	i.sb.WriteByte(')')

	i.sb.WriteString(" VALUES")
	i.args = make([]any, 0, len(fields)*len(i.values))
	for j := 0; j < len(i.values); j++ {
		if j > 0 {
			i.sb.WriteByte(',')
		}
		i.sb.WriteByte('(')
		val := i.values[j]

		//fdVal := reflect.ValueOf(val).Elem()
		refVal := i.valCreator(val, i.model)
		for k, field := range fields {
			if k > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteByte('?')

			//i.args = append(i.args, fdVal.Field(field.Index).Interface())
			fdVal, err := refVal.Field(field.GoName)
			if err != nil {
				return nil, err
			}
			i.args = append(i.args, fdVal)
		}
		i.sb.WriteByte(')')
	}

	// 构造 ON DUPLICATE KEY 部分
	if i.onDuplicate != nil {
		err = i.core.dialect.buildConflictKey(&i.builder, i.onDuplicate)
		if err != nil {
			return nil, err
		}
	}

	i.sb.WriteByte(';')
	return &Query{
		SQL:  i.sb.String(),
		Args: i.args,
	}, err
}
