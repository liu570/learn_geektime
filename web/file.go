package web

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type FileUploader struct {
	// FileField 对应文件在表单中的字段名字
	FileField string

	// DstPathFunc 用于计算目标路径
	DstPathFunc func(fh *multipart.FileHeader) string
}

func (f *FileUploader) Handle() HandleFunc {
	return func(ctx *Context) {
		// FormFile 从表单获取文件数据，需要传入一个 string 参数按名索引
		file, header, err := ctx.Req.FormFile(f.FileField)
		if err != nil {
			ctx.RespStatusCode = 500
			ctx.RespData = []byte("上传失败")
			return
		}
		// 创建目标文件
		dst, err := os.OpenFile(filepath.Join(f.DstPathFunc(header), header.Filename),
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)

		if err != nil {
			ctx.RespStatusCode = 500
			ctx.RespData = []byte("上传失败")
			// 此地可以打日志记录为何上传失败了，但是不能将失败信息返回到前端界面
			return
		}
		// buf为空的话默认 32kb
		// 将 file 文件按 buf 指定的大小，循环复制进入 dst
		io.CopyBuffer(dst, file, nil)
	}
}

type FileDownloader struct {
	Dir string
}

func (f *FileDownloader) Handle() HandleFunc {
	return func(ctx *Context) {
		file, err := ctx.QueryValue("file")
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError // 500
			ctx.RespData = []byte("文件找不到")
			return
		}
		// 因为此地 file 可能是乱七八糟的东西 所以引入 clean 函数 处理一下
		path := filepath.Join(f.Dir, filepath.Clean(file))
		// base 返回文件名的最后一段
		filename := filepath.Base(path)
		header := ctx.Resp.Header()
		// Content-Disposition 指向 attachment 的时候就是表明保存到本地 同时这里还设置的文件名
		header.Set("Content-Disposition", "attachment;filename="+filename)
		header.Set("Content-Description", "File Transfer")
		// Content-Type 值为 octet-stream 的时候代表的是二进制文件 如果知道确切的类型就可以换成别的
		header.Set("Content-Type", "application/octet-stream")
		// binary 参数表示直接传输
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")

		// 文件下载本质上就是将文件写入返回
		http.ServeFile(ctx.Resp, ctx.Req, path)
		//ctx.Resp.Write()
	}
}

type StaticResourceHandlerOption func(handler *StaticResourceHandler)

func WithStaticResourceDir(dir string) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.dir = dir
	}
}

func WithMoreContentType(extMap map[string]string) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		for ext, contentType := range extMap {
			handler.extensionContentTypeMap[ext] = contentType
		}
	}
}

func WithMaxFileSize(fileSize int) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.maxFileSize = fileSize
	}
}

type StaticResourceHandler struct {
	dir                     string
	extensionContentTypeMap map[string]string
	// 我要缓存住静态资源，但不是缓存所有的，并且要控制住内存消耗
	// 第三方库 lru "github.com/hashicorp/golang-lru/v2"
	cache       *lru.Cache[string, *cacheItem]
	maxFileSize int
}

func NewStaticResourceHandler(opts ...StaticResourceHandlerOption) *StaticResourceHandler {
	cache, _ := lru.New[string, *cacheItem](100)
	res := &StaticResourceHandler{
		cache:       cache,
		maxFileSize: 50 * 1024 * 1024,
		extensionContentTypeMap: map[string]string{
			// 这里根据自己的需要不断添加
			// 此地为设置的扩展名
			"jpeg":      "image/jpeg",
			"jpe":       "image/jpeg",
			"jpg":       "image/jpeg",
			"learn_png": "image/learn_png",
			"pdf":       "image/pdf ",
		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

// 假设我们请求 / static/come_on_baby.jpg
func (s *StaticResourceHandler) Handle(ctx *Context) {
	file, ok := ctx.pathParams["file"]
	if !ok {
		ctx.RespStatusCode = 500
		ctx.RespData = []byte("文件不存在")
		return
	}
	header := ctx.Resp.Header()
	itm, ok := s.cache.Get(file)
	if ok {
		header.Set("Content-Type", itm.contentType)
		ctx.RespStatusCode = 200
		ctx.RespData = itm.data
		return
	}
	path := filepath.Join(s.dir, filepath.Clean(file))
	f, err := os.Open(path)
	if err != nil {
		ctx.RespStatusCode = 500
		ctx.RespData = []byte("文件不存在")
		return
	}
	data, err := io.ReadAll(f)
	if err != nil {
		ctx.RespStatusCode = 500
		ctx.RespData = []byte("文件不存在")
		return
	}
	newItm := &cacheItem{
		data:        data,
		contentType: s.extensionContentTypeMap[filepath.Ext(file)],
	}
	// filepath.Ext(file) 返回文件的扩展名
	// 这里需要优化当文件是大型文件就不缓存
	if len(data) <= s.maxFileSize {
		s.cache.Add(file, newItm)
	}
	header.Set("Content-Type", newItm.contentType)
	ctx.RespStatusCode = 200
	ctx.RespData = newItm.data
}

type cacheItem struct {
	data        []byte
	contentType string
}
