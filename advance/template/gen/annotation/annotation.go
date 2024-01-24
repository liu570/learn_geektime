package annotation

import (
	"go/ast"
	"strings"
)

// Annotations 注解类（复数) 用于 AST 模板编程存储一系列注解
// ast.Node（接口） 节点主要有 3 类：表达式和类型节点、语句节点和声明节点。所有节点类型都实现 Node 接口。
type Annotations[N ast.Node] struct {
	Node N
	Ans  []Annotation
}

// Get 使用 key 值查询对应的 value
func (a Annotations[N]) Get(key string) (Annotation, bool) {
	for _, an := range a.Ans {
		if an.Key == key {
			return an, true
		}
	}
	return Annotation{}, false
}

// Annotation 注解类 用于 AST 模板编程存储单个注解键值对
type Annotation struct {
	Key   string
	Value string
}

// newAnnotations 用于创建 Annotations
// 根据 CommentGroup 中的注释来解析对应注解
func newAnnotations[N ast.Node](n N, cg *ast.CommentGroup) Annotations[N] {
	if cg == nil || len(cg.List) == 0 {
		return Annotations[N]{Node: n}
	}
	ans := make([]Annotation, 0, len(cg.List))
	for _, c := range cg.List {
		text, ok := extractContent(c)
		if !ok {
			continue
		}
		if strings.HasPrefix(text, "@") {
			segs := strings.SplitN(text, " ", 2)
			key := segs[0][1:]
			val := ""
			if len(segs) == 2 {
				val = segs[1]
			}
			ans = append(ans, Annotation{
				Key:   key,
				Value: val,
			})
		}
	}
	return Annotations[N]{
		Node: n,
		Ans:  ans,
	}
}

// extractContent 用于提取单个注释 ast.Comment 中的内容
// 即是用于去除注释符号
func extractContent(c *ast.Comment) (string, bool) {
	text := c.Text
	if strings.HasPrefix(text, "// ") {
		return text[3:], true
	} else if strings.HasPrefix(text, "/* ") {
		length := len(text)
		return text[3 : length-2], true
	}
	return "", false
}
