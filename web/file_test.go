package web

import (
	"testing"
)

func TestFileDownloader_Handle(t *testing.T) {
	s := NewHTTPServer()
	s.Get("/download", (&FileDownloader{
		// 下载的文件所在路径
		Dir: "./learn_png/download",
	}).Handle())
	s.Start(":8081")
}

func TestStaticResourceHandler_Handle(t *testing.T) {
	s := NewHTTPServer()
	handler := NewStaticResourceHandler(WithStaticResourceDir("./learn_png/img"))
	s.Get("/img/:file", handler.Handle)
	s.Start(":8081")
}
