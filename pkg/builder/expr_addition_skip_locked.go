package builder

import "github.com/xoctopus/sqlx/pkg/frag"

func SkipLocked() Addition {
	return AsAddition(addition_SKIP_LOCKD, frag.Lit("SKIP LOCKED"))
}
