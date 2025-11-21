package builder

import "github.com/xoctopus/sqlx/pkg/frag"

// ForUpdate pls be careful
func ForUpdate() Addition {
	return AsAddition(addition_FOR_UPDATE, frag.Lit("FOR UPDATE"))
}
