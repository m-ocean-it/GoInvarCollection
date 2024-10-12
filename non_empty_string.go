package invarcol

import invar "github.com/m-ocean-it/GoInvar"

var nonEmptyStringInvariants = []invar.Invariant[string]{
	{
		Name:  "string must be non-empty",
		Check: func(x string) bool { return len(x) > 0 },
	},
}

type NonEmptyString invar.InvariantsHolder[string]

func NewNonEmptyString(s string) NonEmptyString {
	return invar.New(s, nonEmptyStringInvariants)
}

func TryNewNonEmptyString(s string) (NonEmptyString, error) {
	return invar.TryNew(s, nonEmptyStringInvariants)
}
