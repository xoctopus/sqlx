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

	queryMatched bool
	argsMatched  bool
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

	m.queryMatched = m.query == q

	m.argsMatched = true
	if m.checkArgs {
		if len(args) != len(m.args) {
			m.argsMatched = false
		} else if len(args) == 0 {
			m.argsMatched = true
		} else {
			m.argsMatched = reflect.DeepEqual(args, m.args)
		}
	}

	return m.queryMatched && m.argsMatched
}

func (m *fragmentMatcher[A]) NormalizeActual(actual A) any {
	if frag.IsNil(actual) {
		return ""
	}
	q, args := frag.Collect(context.Background(), actual)

	return q + " | " + fmt.Sprintf("%v", args)
}

func (m *fragmentMatcher[A]) NormalizeExpect() any {
	return m.query + " | " + fmt.Sprintf("%v", m.args)
}
