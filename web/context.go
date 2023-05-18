package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type HandleFunc func(ctx *Context)

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	pathParams map[string]string

	// 缓存住你的响应
	// 缓存的响应部分
	// 这部分数据会在最后面刷新
	RespStatusCode int
	// RespData []byte
	RespData []byte

	// 命中路由后将路由路径写入context
	MatchedRoute string

	tplEngine TemplateEngine

	UserValues map[string]any
}

func (ctx *Context) BindJSON(val any) error {
	if ctx.Req.Body == nil {
		return errors.New("web: body 为空")
	}
	// NewDecoder 用于从输入流中读取和解码JSON对象的示例
	decoder := json.NewDecoder(ctx.Req.Body)
	// decoder中的Decode方法是用于 将从输入流中读取的对象转化为 go语言中的对象
	return decoder.Decode(val)
}

func (ctx *Context) FormValue(key string) (string, error) {
	// ParseForm 用于解析 Form 和 PostForm
	err := ctx.Req.ParseForm()
	if err != nil {
		return "", err
	}
	return ctx.Req.FormValue(key), nil
}

func (ctx *Context) QueryValue(key string) (string, error) {
	params := ctx.Req.URL.Query()
	vals, ok := params[key]
	if !ok || len(vals) == 0 {
		return "", nil
	}
	return vals[0], nil
}

func (ctx *Context) PathValue(key string) (string, error) {
	value, ok := ctx.pathParams[key]

	if !ok {
		return "", errors.New("key nod found")
	}
	return value, nil
}

// 同时根据用户的不同需求需要不同类型的数据
// 我们需要搬砖似的补充许多各种类型的方法
// 大明老师认为我们可以抽象出一个返回值结构体并在上方添加类型转换

func (ctx *Context) PathValueV1(key string) StringValue {
	val, ok := ctx.pathParams[key]
	if !ok {
		return StringValue{err: errors.New("key not found")}
	}
	return StringValue{val: val}
}

type StringValue struct {
	val string
	err error
}

func (s StringValue) ToInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}

// / 处理输出

func (ctx *Context) RespJson(code int, val any) error {
	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}
	ctx.RespData = bs
	ctx.RespStatusCode = code
	return err
}

func (ctx *Context) Render(tplName string, val any) error {
	data, err := ctx.tplEngine.Render(ctx.Req.Context(), tplName, val)

	if err != nil {
		return err
	}

	ctx.RespData = data
	ctx.RespStatusCode = http.StatusOK

	return err

}

func (ctx *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(ctx.Resp, cookie)
}
