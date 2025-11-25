package internal

import "github.com/xoctopus/x/contextx"

var (
	CtxTableName  = contextx.NewT[string]()
	CtxTableAlias = contextx.NewT[string]()
)
