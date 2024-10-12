package goinvarcol

import goinvar "github.com/m-ocean-it/GoInvar"

var positiveIntConditions = []goinvar.Condition[int]{
	func(x int) bool { return x > 0 },
}

type PositiveInt goinvar.Invariant[int]

func NewPositiveInt(n int) PositiveInt {
	return goinvar.New(n, positiveIntConditions)
}

func TryNewPositiveInt(n int) (PositiveInt, error) {
	return goinvar.TryNew(n, positiveIntConditions)
}
