package frag

import (
	"bytes"
	"context"
	"database/sql"
	"strings"
	"text/scanner"

	"github.com/xoctopus/x/misc/must"
)

type (
	NamedArg  = sql.NamedArg
	NamedArgs map[string]any
)

func Query(query string, args ...any) Fragment {
	if len(args) == 0 {
		return &pair{query: query}
	}

	frag := &pair{
		query: query,
		args:  make([]any, 0, len(args)),
		set:   make(NamedArgs),
	}

	for _, arg := range args {
		switch x := arg.(type) {
		case NamedArgs:
			for k := range x {
				frag.set[k] = x[k]
			}
		case NamedArg:
			frag.set[x.Name] = x.Value
		default:
			frag.args = append(frag.args, x)
		}
	}

	return frag
}

type pair struct {
	query string
	args  []any
	set   NamedArgs
}

func (p *pair) IsNil() bool { return p.query == "" && len(p.args) == 0 }

func (p *pair) Frag(ctx context.Context) Iter {
	return func(yield func(string, []any) bool) {
		s := &scanner.Scanner{}
		s.Init(strings.NewReader(p.query))
		s.Error = func(s *scanner.Scanner, msg string) {}

		idx := 0
		tmp := bytes.NewBuffer(nil)

		for c := s.Next(); c != scanner.EOF; c = s.Next() {
			switch c {
			case '@':
				if tmp.Len() > 0 {
					if !yield(tmp.String(), nil) {
						return
					}
				}
				tmp.Reset()

				name := bytes.NewBuffer(nil)
				for {
					c = s.Next()
					if c == scanner.EOF {
						break
					}
					if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
						name.WriteRune(c)
						continue
					}
					tmp.WriteRune(c)
					break
				}
				if name.Len() > 0 {
					arg := name.String()
					v, ok := p.set[arg]
					must.BeTrueF(ok, "missing named argument query: %s arg: %s", p.query, arg)
					for query, args := range ArgIter(ctx, v) {
						yield(query, args)
					}
				}
			case '?':
				if tmp.Len() > 0 {
					yield(tmp.String(), nil)
				}
				tmp.Reset()
				must.BeTrueF(
					idx < len(p.args),
					"missing named argument query: %s. given %d arguments at %d",
					p.query, len(p.args), idx,
				)
				arg := p.args[idx]
				for query, args := range ArgIter(ctx, arg) {
					if !yield(query, args) {
						return
					}
				}
				idx++
			default:
				tmp.WriteRune(c)
			}
		}
		if tmp.Len() > 0 {
			if !yield(tmp.String(), nil) {
				return
			}
		}
	}
}
