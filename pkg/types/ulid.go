package types

import (
	"math/big"
	"strings"

	"github.com/oklog/ulid/v2"
)

func MakeULID() ULID {
	return ULID{ULID: ulid.Make()}
}

type ULID struct {
	ulid.ULID
}

func (ULID) DBType(driver string) string {
	switch strings.ToLower(driver) {
	case "postgres", "duckdb":
		return "uuid"
	case "sqlite":
		return "blob"
	default:
		return "binary"
	}
}

func (u ULID) DBFixedWidth(driver string) *uint64 {
	if u.DBType(driver) == "binary" {
		return new(uint64(16))
	}
	return nil
}

func (u ULID) BigInt() *big.Int {
	i := new(big.Int).SetBytes(u.ULID[:])
	return i
}
