package frag

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"sync"

	"github.com/xoctopus/x/iterx"
)

// Alias generates a deterministic alias name for a given table and column.
// If the combined length is shorter than 64 characters (minus separators),
// it returns a simple "table__column" alias. Otherwise, it generates a hashed
// form using a cached SHA1-based abbreviation to keep the alias short and stable.
func Alias(table, column string) string {
	if len(table)+len(column) < 64-2 {
		return table + "__" + column
	}
	return a.hash(table) + "__" + a.hash(column)
}

type aliases struct{ m sync.Map }

var a = &aliases{}

func (as *aliases) hash(name string) string {
	v, _ := as.m.LoadOrStore(name, sync.OnceValue(func() string {
		hash := fmt.Sprintf("%x", sha1.Sum([]byte(name)))
		parts := strings.Split(name, "_")

		return strings.Join(
			iterx.Values(
				iterx.MapSlice(parts, func(s string) string {
					if len(s) >= 1 {
						return strings.ToUpper(s[0:1])
					}
					return ""
				}),
			),
			"",
		) + "_" + hash[0:8]
	}))

	return v.(func() string)()
}
