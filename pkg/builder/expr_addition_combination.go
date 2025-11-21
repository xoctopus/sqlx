package builder

import (
	"context"

	"github.com/xoctopus/sqlx/pkg/frag"
)

type CombinationAddition interface {
	Addition

	All(SelectStatement) CombinationAddition
	Distinct(SelectStatement) CombinationAddition
}

func Union() CombinationAddition {
	return &combination{operator: "UNION"}
}

func Expect() CombinationAddition {
	return &combination{operator: "EXCEPT"}
}

func Intersect() CombinationAddition {
	return &combination{operator: "INTERSECT"}
}

type combination struct {
	operator string // operator UNION | INTERSECT | EXCEPT
	method   string // method ALL | DISTINCT
	stmt     SelectStatement
}

func (combination) Type() AdditionType {
	return addition_COMBINATION
}

func (c combination) All(stmt SelectStatement) CombinationAddition {
	c.method = "ALL"
	c.stmt = stmt
	return &c
}

func (c combination) Distinct(stmt SelectStatement) CombinationAddition {
	c.method = "DISTINCT"
	c.stmt = stmt
	return &c
}

func (c *combination) IsNil() bool {
	return c == nil || frag.IsNil(c.stmt)
}

func (c *combination) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		yield(c.operator+" ", nil)
		if c.method != "" {
			yield(c.method+" ", nil)
		}
		for q, args := range c.stmt.Frag(ctx) {
			yield(q, args)
		}
	}
}
