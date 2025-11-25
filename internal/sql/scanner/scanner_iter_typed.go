package scanner

import (
	"context"
	"iter"
)

type TScanner[T any] struct {
	next func(v *T) error
}

func (*TScanner[T]) New() any {
	return new(T)
}

func (t *TScanner[T]) Next(v any) error {
	return t.next(v.(*T))
}

func Recv[T any](next func(v *T) error) ScanIter {
	return &TScanner[T]{next: next}
}

type RecvFunc[M any] func(ctx context.Context, recv func(v *M) error) error

func (f RecvFunc[M]) Items(ctx context.Context) iter.Seq2[*M, error] {
	return func(yield func(item *M, err error) bool) {
		var cancel context.CancelFunc

		ctx, cancel = context.WithCancel(ctx)
		defer cancel()

		items := make(chan *M)
		errch := make(chan error)
		go func() {
			defer close(items)
			defer close(errch)

			errch <- f(
				ctx,
				func(item *M) error {
					items <- item
					return nil
				},
			)
		}()

		for {
			select {
			case item := <-items:
				if !yield(item, nil) {
					cancel()
					<-errch
					return
				}
				continue
			case err := <-errch:
				if err != nil {
					if !yield(nil, err) {
						return
					}
				}
				return
			}
		}
	}
}
