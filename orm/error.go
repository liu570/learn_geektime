package orm

import "learn_geektime/orm/internal/errs"

// 将内部的 internal error 暴露出去
var (
	// ErrNoRows 代表没有找到数据
	ErrNoRows = errs.ErrNoRows
)
