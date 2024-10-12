package invarcol

import (
	"errors"

	invar "github.com/m-ocean-it/GoInvar"
)

type NonNilPointer[T any] invar.InvariantsHolder[*T]

func getNonNilPointerInvariants[T any]() []invar.Invariant[*T] {
	return []invar.Invariant[*T]{
		func(ptr *T) error {
			if ptr != nil {
				return nil
			}

			return errors.New("pointer must not be nil")
		},
	}
}

func NewNonNilPointer[T any](ptr *T) NonNilPointer[T] {
	return invar.New(ptr, getNonNilPointerInvariants[T]())
}

func TryNewNonNilPointer[T any](ptr *T) (NonNilPointer[T], error) {
	return invar.TryNew(ptr, getNonNilPointerInvariants[T]())
}
