package goinvarcol

import goinvar "github.com/m-ocean-it/GoInvar"

var nonEmptyStringConditions = []goinvar.Condition[string]{
	func(x string) bool { return len(x) > 0 },
}

type NonEmptyString goinvar.Invariant[string]

func NewNonEmptyString(s string) NonEmptyString {
	return goinvar.New(s, nonEmptyStringConditions)
}

func TryNewNonEmptyString(s string) (NonEmptyString, error) {
	return goinvar.TryNew(s, nonEmptyStringConditions)
}
