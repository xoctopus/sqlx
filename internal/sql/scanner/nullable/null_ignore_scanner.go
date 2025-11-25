package nullable

import (
	"database/sql"
	_ "unsafe"
)

//go:linkname convertAssign database/sql.convertAssign
func convertAssign(dst, src any) error

func NewNullIgnoreScanner(dst any) *NullIgnoreScanner {
	return &NullIgnoreScanner{
		dst: dst,
	}
}

type NullIgnoreScanner struct {
	dst any
}

func (s *NullIgnoreScanner) Scan(src any) error {
	if scanner, ok := s.dst.(sql.Scanner); ok {
		return scanner.Scan(src)
	}
	if src == nil {
		return nil
	}
	return convertAssign(s.dst, src)
}

type Empty struct{}

func (scanner *Empty) Scan(any) error {
	return nil
}
