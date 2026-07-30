package adaptor

import (
	"net/url"
	"testing"

	. "github.com/xoctopus/x/testx"
)

func TestDatabaseNameFromDSN(t *testing.T) {
	parse := func(raw string) *url.URL {
		u, err := url.Parse(raw)
		Expect(t, err, Succeed())
		return u
	}

	t.Run("WithLeadingSlash", func(t *testing.T) {
		Expect(t, DatabaseNameFromDSN(parse("mysql://root@localhost:3306/mydb")), Equal("mydb"))
	})

	t.Run("NestedPath", func(t *testing.T) {
		Expect(t, DatabaseNameFromDSN(parse("mysql://root@localhost:3306/a/b")), Equal("a/b"))
	})

	t.Run("EmptyPath", func(t *testing.T) {
		Expect(t, DatabaseNameFromDSN(parse("mysql://root@localhost:3306")), Equal(""))
	})

	t.Run("RootOnly", func(t *testing.T) {
		Expect(t, DatabaseNameFromDSN(parse("mysql://root@localhost:3306/")), Equal(""))
	})

	t.Run("WithQuery", func(t *testing.T) {
		Expect(
			t,
			DatabaseNameFromDSN(parse("postgres://u@h:5432/app?sslmode=disable")),
			Equal("app"),
		)
	})

	t.Run("DirectPath", func(t *testing.T) {
		Expect(t, DatabaseNameFromDSN(&url.URL{Path: "/dbname"}), Equal("dbname"))
		Expect(t, DatabaseNameFromDSN(&url.URL{Path: "dbname"}), Equal("dbname"))
		Expect(t, DatabaseNameFromDSN(&url.URL{}), Equal(""))
	})
}
