package builder

import (
	"strings"

	"github.com/xoctopus/sqlx/pkg/frag"
)

// ForShare builds a FOR SHARE lock clause.
func ForShare(modifiers ...string) Addition {
	return AsAddition(
		addition_FOR_SHARE,
		frag.Lit(strings.Join(append([]string{"FOR SHARE"}, modifiers...), " ")),
	)
}

// ForShareSkipLocked builds FOR SHARE SKIP LOCKED.
func ForShareSkipLocked() Addition {
	return ForShare("SKIP LOCKED")
}

// ForShareNoWait builds FOR SHARE NOWAIT.
func ForShareNoWait() Addition {
	return ForShare("NOWAIT")
}
