package builder

import (
	"context"
	"strings"

	"github.com/xoctopus/sqlx/pkg/frag"
)

func Comment(lines ...string) Addition {
	c := &comment{}
	for _, text := range lines {
		for _, line := range strings.Split(text, "\n") {
			line = strings.TrimSpace(line)
			if len(line) > 0 {
				c.lines = append(c.lines, line)
			}
		}

	}
	return c
}

type comment struct {
	lines []string
}

func (c *comment) Type() AdditionType {
	return addition_COMMENT
}

func (c *comment) IsNil() bool {
	return c == nil || len(c.lines) == 0
}

func (c *comment) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		for i, line := range c.lines {
			if i > 0 {
				yield("\n", nil)
			}
			yield("-- "+line, nil)
		}
	}
}
