package testutil

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/frag"
)

func BeFragmentForQuery(query string, args ...any) testx.Matcher[frag.Fragment] {
	return &fragmentMatcher[frag.Fragment]{
		query: strings.TrimSpace(query),
		args:  args,
	}
}

func BeFragment(query string, args ...any) testx.Matcher[frag.Fragment] {
	return &fragmentMatcher[frag.Fragment]{
		query:     strings.TrimSpace(query),
		args:      args,
		checkArgs: true,
	}
}

type fragmentMatcher[A frag.Fragment] struct {
	query     string
	args      []any
	checkArgs bool

	queryNotEqual bool
	argsNotEqual  bool
}

func (m *fragmentMatcher[A]) Negative() bool {
	return false
}

func (m *fragmentMatcher[A]) Action() string {
	return "Be Frag"
}

func (m *fragmentMatcher[A]) Match(actual A) bool {
	if frag.IsNil(actual) {
		return m.query == ""
	}
	q, args := frag.Collect(context.Background(), actual)
	if len(m.args) == 0 && len(args) == 0 {
		m.queryNotEqual = !(m.query == q)
		return !m.queryNotEqual
	}

	if m.query == q {
		return true
	}

	m.queryNotEqual = true

	if m.checkArgs {
		return reflect.DeepEqual(m.args, args)
	}

	m.argsNotEqual = true

	return false
}

func (m *fragmentMatcher[A]) NormalizeActual(actual A) any {
	if frag.IsNil(actual) {
		return ""
	}
	q, args := frag.Collect(context.Background(), actual)

	if m.queryNotEqual && m.argsNotEqual {
		return fmt.Sprintf("%s | %v", q, args)
	}

	if m.queryNotEqual {
		return q
	}

	return fmt.Sprintf("%v", args)
}

func (m *fragmentMatcher[A]) NormalizeExpect() any {
	if m.queryNotEqual && m.argsNotEqual {
		return fmt.Sprintf("%s | %v", m.query, m.args)
	}

	if m.queryNotEqual {
		return m.query
	}

	return fmt.Sprintf("%v", m.args)
}
