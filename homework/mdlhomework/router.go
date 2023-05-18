package mdlhomework

import (
	"strings"
)

type router struct {
	trees map[string]*node
}

func (r *router) addRoute(method string, path string, handleFunc HandleFunc) {
	if path == "" {
		panic("路径为空")
	}
	root, ok := r.trees[method]
	if !ok {
		root = &node{path: "/"}
		r.trees[method] = root
	}
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		if seg == "" {
			continue
		}
		root = root.childOrCreate(seg)
	}
	root.handle = handleFunc
}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	if path == "" {
		panic("路径为空")
	}
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	pathParams := make(map[string]string)
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		if seg == "" {
			continue
		}
		root = root.findChild(seg)
		if root != nil && root.path[0] == ':' {
			pathParams[root.path[1:]] = seg
		}
	}
	if root == nil {
		return nil, false
	}
	return &matchInfo{node: root, pathParams: pathParams}, root.handle != nil
}

func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

type node struct {
	path     string
	children map[string]*node
	handle   HandleFunc

	starChild  *node
	paramChild *node
}

func (n *node) childOrCreate(path string) *node {
	// 此地不支持连续两个通配符 eg：/user/*/*
	if path == "*" {
		if n.starChild == nil {
			if n.paramChild != nil {
				panic("Not allow both parameter routes and wildcard " +
					"routes to be registered under the same path")
			}
			n.starChild = &node{path: path}
		}
		return n.starChild
	}
	if path[0] == ':' {
		if n.paramChild == nil {
			if n.starChild != nil {
				panic("Not allow both parameter routes and wildcard " +
					"routes to be registered under the same path")
			}
			n.paramChild = &node{path: path}
		}
		return n.paramChild
	}

	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[path]
	if !ok {
		child = &node{path: path}
		n.children[path] = child
	}
	return child
}

func (n *node) findChild(path string) *node {
	if n == nil {
		return nil
	}
	if n.children == nil {
		if n.starChild != nil {
			return n.starChild
		} else if n.paramChild != nil {
			return n.paramChild
		}
		return nil
	}
	child, ok := n.children[path]
	if !ok {
		if n.starChild != nil {
			return n.starChild
		} else if n.paramChild != nil {
			return n.paramChild
		}
		return nil
	}
	return child
}

type matchInfo struct {
	node       *node
	pathParams map[string]string
}
