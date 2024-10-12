package invarcol

import invar "github.com/m-ocean-it/GoInvar"

var nonEmptyStringConditions = []invar.Condition[string]{
	func(x string) bool { return len(x) > 0 },
}

type NonEmptyString invar.Invariant[string]

func NewNonEmptyString(s string) NonEmptyString {
	return invar.New(s, nonEmptyStringConditions)
}

func TryNewNonEmptyString(s string) (NonEmptyString, error) {
	return invar.TryNew(s, nonEmptyStringConditions)
}
