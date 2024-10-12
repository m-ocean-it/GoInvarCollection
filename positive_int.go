package invarcol

import (
	"errors"

	invar "github.com/m-ocean-it/GoInvar"
)

var positiveIntInvariants = []invar.Invariant[int]{
	func(x int) error {
		if x > 0 {
			return nil
		}

		return errors.New("number must be positive")
	},
}

type PositiveInt invar.InvariantsHolder[int]

func NewPositiveInt(n int) PositiveInt {
	return invar.New(n, positiveIntInvariants)
}

func TryNewPositiveInt(n int) (PositiveInt, error) {
	return invar.TryNew(n, positiveIntInvariants)
}
