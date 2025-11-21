package builder

import (
	"context"

	"github.com/xoctopus/sqlx/pkg/frag"
)

type ConditionMarker interface {
	asCondition()
}

func And(conditions ...frag.Fragment) SqlCondition {
	return compose("AND", process(conditions))
}

func Or(conditions ...frag.Fragment) SqlCondition {
	return compose("OR", process(conditions))
}

func Xor(conditions ...frag.Fragment) SqlCondition {
	return compose("XOR", process(conditions))
}

type SqlCondition interface {
	frag.Fragment

	ConditionMarker
}

func AsCond(f frag.Fragment) *Condition {
	switch x := f.(type) {
	case *Condition:
		return x
	default:
		return &Condition{expr: x}
	}
}

type Condition struct {
	expr frag.Fragment

	ConditionMarker
}

func (c *Condition) IsNil() bool {
	return c == nil || frag.IsNil(c.expr)
}

func (c *Condition) Frag(ctx context.Context) frag.Iter {
	if frag.IsNil(c.expr) {
		return nil
	}
	return c.expr.Frag(ctx)
}

// compose conditions with operator
func compose(operator string, conditions []SqlCondition) SqlCondition {
	return &ComposedCondition{
		operator:   operator,
		conditions: conditions,
	}
}

// process conditions, filter nil and where expr
func process(conditions []frag.Fragment) []SqlCondition {
	final := make([]SqlCondition, 0, len(conditions))

	for i := range conditions {
		c := conditions[i]
		if w, ok := c.(*where); ok {
			c = w.condition
		}
		if frag.IsNil(c) {
			continue
		}
		final = append(final, AsCond(c))
	}
	return final
}

type ComposedCondition struct {
	ConditionMarker

	operator   string
	conditions []SqlCondition
}

func (c *ComposedCondition) IsNil() bool {
	if c == nil || c.operator == "" || len(c.conditions) == 0 {
		return true
	}
	for i := range c.conditions {
		if !frag.IsNil(c.conditions[i]) {
			return false
		}
	}
	return true
}

func (c *ComposedCondition) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		if len(c.conditions) == 1 {
			for q, args := range c.conditions[0].Frag(ctx) {
				yield(q, args)
			}
			return
		}
		for i, cond := range c.conditions {
			if i > 0 {
				yield(" "+c.operator+" ", nil)
			}
			for q, args := range frag.Block(cond).Frag(ctx) {
				yield(q, args)
			}
		}
	}
}
