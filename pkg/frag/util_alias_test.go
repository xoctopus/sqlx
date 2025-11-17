package frag_test

import (
	"crypto/sha1"
	"fmt"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/frag"
)

func TestAlias(t *testing.T) {
	Expect(t, frag.Alias("table", "column"), Equal("table__column"))

	tab := "too_________________long___________table_______________name"
	col := "too_________________long___________column______________name"
	hashTab := fmt.Sprintf("%x", sha1.Sum([]byte(tab)))
	hashCol := fmt.Sprintf("%x", sha1.Sum([]byte(col)))

	Expect(t, frag.Alias(tab, col), Equal("TLTN_"+hashTab[0:8]+"__TLCN_"+hashCol[0:8]))
}
