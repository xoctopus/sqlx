package session

import (
	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/internal/sql/scanner"
)

var (
	Scan = scanner.Scan
	Open = adaptor.Open
)
