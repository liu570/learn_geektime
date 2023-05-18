package web

import (
	"bytes"
	"context"
	"html/template"
	"io/fs"
)

type TemplateEngine interface {
	// 第一个返回值：渲染好的数据
	Render(ctx context.Context, tplName string, data any) ([]byte, error)
}

type GoTemplateEngine struct {
	T *template.Template
	// 也可以考虑设计为map[string]*template.Template
	// 但是其实没太大必要，因为 template.Template 本身就提供了按名索引的功能
}

func (g *GoTemplateEngine) Render(ctx context.Context, tplname string, data any) ([]byte, error) {
	res := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(res, tplname, data)
	return res.Bytes(), err
}

// 以下三个方法(为管理模板的方法)，可以加也可以不加，并不是 web 框架的核心功能完全可以由用户自己实现

func (g *GoTemplateEngine) LoadFromGlob(pattern string) error {
	var err error
	g.T, err = template.ParseGlob(pattern)
	return err
}
func (g *GoTemplateEngine) LoadFromFiles(filenames ...string) error {
	var err error
	g.T, err = template.ParseFiles(filenames...)
	return err
}
func (g *GoTemplateEngine) LoadFromFS(fs fs.FS, patterns ...string) error {
	var err error
	g.T, err = template.ParseFS(fs, patterns...)
	return err
}
