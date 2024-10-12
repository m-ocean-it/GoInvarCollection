package invarcol

import invar "github.com/m-ocean-it/GoInvar"

var positiveIntConditions = []invar.Condition[int]{
	func(x int) bool { return x > 0 },
}

type PositiveInt invar.Invariant[int]

func NewPositiveInt(n int) PositiveInt {
	return invar.New(n, positiveIntConditions)
}

func TryNewPositiveInt(n int) (PositiveInt, error) {
	return invar.TryNew(n, positiveIntConditions)
}
