package scanner

import (
	"reflect"

	"github.com/xoctopus/x/reflectx"
)

type ScanIter interface {
	// New a ptr value for scanning
	New() any
	// Next for receive scanned value
	Next(v any) error
}

type WithColumnReceivers interface {
	ColumnReceivers() map[string]any
}

type ScanIterT[T any] interface {
	New() T
	Next(T) error
}

func ScanIterFor(v any) (ScanIter, error) {
	switch x := v.(type) {
	case ScanIter:
		return x, nil
	default:
		t := reflectx.Deref(reflect.TypeOf(v))
		if t.Kind() == reflect.Slice && t.Elem().Kind() != reflect.Uint8 {
			return &SliceScanIter{
				t: t.Elem(),
				v: reflectx.Indirect(reflect.ValueOf(v)),
			}, nil
		}
		return &SingleScanIter{v: v}, nil
	}
}

type SliceScanIter struct {
	t reflect.Type
	v reflect.Value
}

func (s *SliceScanIter) New() any {
	return reflectx.New(s.t).Addr().Interface()
}

func (s *SliceScanIter) Next(v any) error {
	s.v.Set(reflect.Append(s.v, reflect.ValueOf(v).Elem()))
	return nil
}

type SingleScanIter struct {
	v          any
	hasResults bool
}

func (s *SingleScanIter) New() any {
	return s.v
}

func (s *SingleScanIter) Next(v any) error {
	s.hasResults = true
	return nil
}

func (s *SingleScanIter) MustHasRecord() bool {
	return s.hasResults
}
