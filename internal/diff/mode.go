package diff

import (
	"github.com/xoctopus/x/contextx"
	"github.com/xoctopus/x/flagx"
)

type ActType int8

const (
	ACT_DROP_IDX ActType = iota + 1
	ACT_DROP_COL
	ACT_KEEP_COL
	ACT_RENAME_COL
	ACT_MODIFY_COL
	ACT_CREATE_COL
	ACT_CREATE_IDX
	ACT_CREATE_TABLE
)

type Mode uint8

const (
	MODE_CREATE_TABLE Mode = 1 << iota
	MODE_DRY_RUN
)

var CtxMode = contextx.NewT[flagx.Flagger[Mode]]()
