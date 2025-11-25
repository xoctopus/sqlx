package errors

import (
	"errors"
	"fmt"

	"github.com/xoctopus/errx/pkg/codex"
)

type code int8

const (
	NOTFOUND code = iota + 1
	CONFLICT
	ROLLBACK
)

func (c code) Message() string {
	prefix := fmt.Sprintf("[SQLERROR:%d] ", c)

	switch c {
	case NOTFOUND:
		return prefix + "NOTFOUND"
	case CONFLICT:
		return prefix + "CONFLICT"
	case ROLLBACK:
		return prefix + "ROLLBACK"
	default:
		return prefix + "UNKNOWN"
	}
}

func IsErrNotFound(err error) bool {
	return errors.Is(err, codex.New(NOTFOUND))
}

func IsErrConflict(err error) bool {
	return errors.Is(err, codex.New(NOTFOUND))
}

func IsErrRollback(err error) bool {
	return errors.Is(err, codex.New(ROLLBACK))
}
