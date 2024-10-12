package invarcol

import (
	"errors"

	invar "github.com/m-ocean-it/GoInvar"
)

var nonEmptyStringInvariants = []invar.Invariant[string]{
	func(x string) error {
		if len(x) > 0 {
			return nil
		}

		return errors.New("string must be non-empty")
	},
}

type NonEmptyString invar.InvariantsHolder[string]

func NewNonEmptyString(s string) NonEmptyString {
	return invar.New(s, nonEmptyStringInvariants)
}

func TryNewNonEmptyString(s string) (NonEmptyString, error) {
	return invar.TryNew(s, nonEmptyStringInvariants)
}
