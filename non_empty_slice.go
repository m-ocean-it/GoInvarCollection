package invarcol

import invar "github.com/m-ocean-it/GoInvar"

type NonEmptySlice[T any] invar.InvariantsHolder[[]T]

func NewNonEmptySlice[T any](s []T) NonEmptySlice[T] {
	return invar.New(s, getNonEmptySliceInvariants[T]())
}

func TryNewNonEmptySlice[T any](s []T) (NonEmptySlice[T], error) {
	return invar.TryNew(s, getNonEmptySliceInvariants[T]())
}

func getNonEmptySliceInvariants[T any]() []invar.Invariant[[]T] {
	return []invar.Invariant[[]T]{
		{
			Name:  "slice must be non-empty",
			Check: func(s []T) bool { return len(s) > 0 },
		},
	}
}
