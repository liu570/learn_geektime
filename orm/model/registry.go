package model

import (
	"learn_geektime/orm/internal/errs"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

type Registry interface {
	Get(val any) (*Model, error)
	Register(val any, opts ...ModelOpt) (*Model, error)
}

// registry 该结构体主要用来缓存 models 即是缓存元数据
type registry struct {
	// 使用 sync.Map 是为了线程安全
	models sync.Map
}

func NewRegistry() Registry {
	return &registry{}
}

func (r *registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}
	// Register 只在这里使用 其实我们可以直接传入 typ 进去 可以少调用一次 ValueOf
	return r.Register(val)
}

func (r *registry) Register(val any, opts ...ModelOpt) (*Model, error) {
	res, err := r.parseModel(val)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		err := opt(res)
		if err != nil {
			return nil, err
		}
	}
	typ := reflect.TypeOf(val)
	r.models.Store(typ, res)
	return res, nil
}

// parseModel 支持从标签中提取自定义设置
// 标签形式 orm:"key1=value1,key2=value2"
func (r *registry) parseModel(val any) (*Model, error) {
	if val == nil {
		return nil, errs.ErrInputNil
	}
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointOnly
	}
	typ = typ.Elem()

	numField := typ.NumField()
	fieldMap := make(map[string]*Field, numField)
	colMap := make(map[string]*Field, numField)
	columns := make([]*Field, numField)
	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		ormTagKvs, err := r.parseTag(fd.Tag)
		if err != nil {
			return nil, err
		}
		colName, ok := ormTagKvs[tagKeyColumn]
		if !ok || colName == "" {
			colName = CamelCaseToSnakeCase(fd.Name)
		}
		fdMeta := &Field{
			ColName: colName,
			Type:    fd.Type,
			GoName:  fd.Name,
			Offset:  fd.Offset,
			Index:   i,
		}
		fieldMap[fd.Name] = fdMeta
		colMap[colName] = fdMeta
		columns[i] = fdMeta
	}

	var tableName string
	if tn, ok := val.(TableName); ok {
		tableName = tn.TableName()
	}

	if tableName == "" {
		tableName = CamelCaseToSnakeCase(typ.Name())
	}
	res := &Model{
		TableName: tableName,
		FieldMap:  fieldMap,
		ColumnMap: colMap,
		Fields:    columns,
	}

	return res, nil
}

// column => id
func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag := tag.Get("orm")
	if ormTag == "" {
		// 传入一个空的 map 上层调用就不需要判断 map 是否为空值了
		return map[string]string{}, nil
	}
	kvs := strings.Split(ormTag, ",")
	res := make(map[string]string, len(kvs))
	for _, kv := range kvs {
		val := strings.SplitN(kv, "=", 2)
		if len(val) != 2 {
			return nil, errs.NewErrInvalidTagContent(kv)
		}
		res[val[0]] = val[1]
	}
	return res, nil
}

func CamelCaseToSnakeCase(s string) string {
	var res strings.Builder
	for i, c := range s {
		if unicode.IsUpper(c) {
			if i > 0 && (unicode.IsLower(rune(s[i-1])) || (i+1 < len(s) && unicode.IsLower(rune(s[i+1])))) {
				res.WriteByte('_')
			}
			res.WriteRune(unicode.ToLower(c))
		} else {
			res.WriteRune(c)
		}
	}
	return res.String()
}

// camelCaseToSnakeCase 正则表达式写法
//func camelCaseToSnakeCase(s string) string {
//	re := regexp.MustCompile(`(?<!^)(?=[A-Z])`)
//	return strings.ToLower(re.ReplaceAllString(s, "_"))
//}
