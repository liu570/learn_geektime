package web

import (
	"fmt"
	"strings"
)

// 大明建议，有明确的意义的东西最好可以抽离出来创建一个结构体

// 代表路由
type router struct {
	// trees 代表的森林，HTTP method => 树的根节点
	trees map[string]*node
}

func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

// addRoute 注册路由。
// method 是 HTTP 方法
// - 已经注册了的路由，无法被覆盖。例如 /user/home 注册两次，会冲突
// - path 必须以 / 开始
func (r *router) addRoute(method string, path string, handleFunc HandleFunc) {
	if path == "" {
		panic("web:路由为空字符串")
	}
	if path[0] != '/' {
		panic("web:路由必须以/开头")
	}

	root, ok := r.trees[method]
	if !ok {
		// 表明根节点还没创建
		root = &node{path: "/"}
		r.trees[method] = root
	}
	// 支持静态路由匹配
	cur := root
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		if seg == "" {
			continue
		}
		cur = cur.childOrCreate(seg)
	}
	if cur.handler != nil {
		panic(fmt.Sprintf("web: 路由冲突[%s]", path))
	}
	//绑定相应的 handlerFunc
	cur.handler = handleFunc
	cur.route = path
}

// findRoute 查找路由
func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return &matchInfo{n: root}, true
	}
	cur := root
	segs := strings.Split(path, "/")

	pathParams := make(map[string]string)

	for _, seg := range segs {
		if seg == "" {
			continue
		}
		cur = cur.findChild(seg)
		if cur != nil && cur.path[0] == ':' {
			pathParams[cur.path[1:]] = seg
		}
	}
	if cur == nil {
		return nil, false
	}

	return &matchInfo{n: cur, pathParams: pathParams}, cur.handler != nil
}

type node struct {
	// 维护当前这一段 /a/b/c 中的一段
	path    string
	route   string
	handler HandleFunc
	// 静态路由
	children map[string]*node
	// 通配符匹配
	starChild *node
	// 参数匹配
	paramChild *node
}

func (n *node) childOrCreate(path string) *node {
	// 通配符  /a/*/b
	if path == "*" {
		if n.starChild == nil {
			n.starChild = &node{
				path: path,
			}
		}
		return n.starChild
	}
	// 参数   /a/:bcd
	if path[0] == ':' {
		if n.paramChild == nil {
			n.paramChild = &node{
				path: path,
			}
		}
		return n.paramChild
	}
	// 如果没有子节点则创建子节点
	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[path]
	// 如果子节点没有该项目
	if !ok {
		child = &node{path: path}
		//将创建的孩子节点加入树中
		n.children[path] = child
	}
	// 返回创建的孩子节点
	return child
}

func (n *node) findChild(path string) *node {
	if n == nil {
		return n
	}
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild
		} else if n.starChild != nil {
			return n.starChild
		} else {
			return nil
		}
	}
	child, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild
		} else if n.starChild != nil {
			return n.starChild
		}
	}
	return child
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}
