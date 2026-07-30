package builder

import (
	"strings"

	"github.com/xoctopus/sqlx/pkg/frag"
)

// ForUpdate builds a FOR UPDATE lock clause.
func ForUpdate(modifiers ...string) Addition {
	return AsAddition(
		addition_FOR_UPDATE,
		frag.Lit(strings.Join(append([]string{"FOR UPDATE"}, modifiers...), " ")),
	)
}

// ForUpdateSkipLocked builds FOR UPDATE SKIP LOCKED.
func ForUpdateSkipLocked() Addition {
	return ForUpdate("SKIP LOCKED")
}

// ForUpdateNoWait builds FOR UPDATE NOWAIT.
func ForUpdateNoWait() Addition {
	return ForUpdate("NOWAIT")
}
