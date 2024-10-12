package invarcol

import invar "github.com/m-ocean-it/GoInvar"

var positiveIntInvariants = []invar.Invariant[int]{
	{
		Name:  "number must be positive",
		Check: func(x int) bool { return x > 0 },
	},
}

type PositiveInt invar.InvariantsHolder[int]

func NewPositiveInt(n int) PositiveInt {
	return invar.New(n, positiveIntInvariants)
}

func TryNewPositiveInt(n int) (PositiveInt, error) {
	return invar.TryNew(n, positiveIntInvariants)
}
