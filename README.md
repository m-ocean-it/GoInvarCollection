# Collection of [GoInvar](https://github.com/m-ocean-it/GoInvar) types

## Implemented types

- [Positive integer](https://github.com/m-ocean-it/GoInvarCollection/blob/main/positive_int.go)
- [Non-empty string](https://github.com/m-ocean-it/GoInvarCollection/blob/main/non_empty_string.go)
- [Non-empty slice](https://github.com/m-ocean-it/GoInvarCollection/blob/main/non_empty_slice.go)
- [Non-nil pointer](https://github.com/m-ocean-it/GoInvarCollection/blob/main/non_nil_pointer.go)

### TODO
- [ ] non-empty map
- [ ] positive int64
- [ ] positive int32
- [ ] positive float64
- [ ] positive float32
- [ ] non-nil error
- [ ] string containing non-whitespace characters
- [ ] uuid string
- [ ] email string
- [ ] url string

## Example usage

### Primitive value invariant

```go
package main

import (
    "fmt"

    invar "github.com/m-ocean-it/GoInvar"
    invarcol "github.com/m-ocean-it/GoInvarCollection"
)

func main() {
    positiveInt, err := invarcol.TryNewPositiveInt(2)
    if err != nil {
        panic(err)
    }

    // Here, it's safe to call Unwrap, instead of TryUnwrap, since initialization didn't error above.
    n := invar.Unwrap(positiveInt)

    fmt.Println("positive int is " + fmt.Sprint(n))
}
```

### Encoding struct invariants

#### Use `Invariant[T]`, where `T` is the struct

`person/person.go`:
```go
package person

import (
	invar "github.com/m-ocean-it/GoInvar"
	invarcol "github.com/m-ocean-it/GoInvarCollection"
)

// The struct itself must be private so that it could only be created via the constructor.
type person struct {
	Name       invarcol.NonEmptyString
	Age        invarcol.PositiveInt
	PlacesBeen invarcol.NonEmptySlice[invarcol.NonEmptyString]
}

// ValidPerson is our struct invariant. As an interface, it cannot be directly initialized.
// Also, since the person struct is private, no other package would be able implement that interface.
// The underlying person struct will be accessible via the Invariant.Get method.
type ValidPerson invar.InvariantsHolder[person]

// New is a custom constructor that checks individual field invariants and returns ValidPerson.
// It's also possible to check inter-field invariants within a constructor.
func New(
	name invarcol.NonEmptyString,
	age invarcol.PositiveInt,
	placesBeen invarcol.NonEmptySlice[invarcol.NonEmptyString],
) (ValidPerson, error) {
	p := person{
		Name:       name,
		Age:        age,
		PlacesBeen: placesBeen,
	}

	// It's important to specify non-nilness of those fields as the struct's invariants,
	// since, unfortunately, any interface can be nil in Go. Failing to do so will lead to
	// an error upon TryUnwrap'ing one of those values or a panic upon calling Unwrap.

	return invar.TryNew(p, []invar.Invariant[person]{
		{
			Name:  "name must be initialized",
			Check: func(p person) bool { return p.Name != nil },
		},
		{
			Name:  "age must be initialized",
			Check: func(p person) bool { return p.Age != nil },
		},
		{
			Name:  "placesBeen must be initialized",
			Check: func(p person) bool { return p.PlacesBeen != nil },
		},
	})
}
```

`main.go`:
```go
package main

import (
	"app/person"
	"fmt"

	invar "github.com/m-ocean-it/GoInvar"
	invarcol "github.com/m-ocean-it/GoInvarCollection"
)

func main() {
	// It's okay to call New... instead of TryNew... when you are sure the invariants hold up. It won't panic.
	nonEmptyName := invarcol.NewNonEmptyString("John Doe")
	positiveAge := invarcol.NewPositiveInt(42)

	sliceOfPlaces := []invarcol.NonEmptyString{
		invarcol.NewNonEmptyString("London"),
	}

	placesBeen := invarcol.NewNonEmptySlice(sliceOfPlaces)

	// sliceOfPlaces[0] = nil // <--------------------------------------- TRY UNCOMMENTING

	p, err := person.New(nonEmptyName, positiveAge, placesBeen)
	if err != nil {
		panic(err)
	}

	unwrappedPerson, err := invar.TryUnwrap(p)
	if err != nil {
		panic(err)
	}

	// Here, it's safe to use Unwrap instead of TryUnwrap, because string and
	// int are fully copied when constructing an InvariantsHolders, therefore
	// there are no external pointers to those values.
	//
	// (But you can still stay on the safe side and use TryUnwrap, if you feel like it.)
	unwrappedName := invar.Unwrap(unwrappedPerson.Name)
	unwrappedAge := invar.Unwrap(unwrappedPerson.Age)

	// Here, it's better to call TryUnwrap, since slice could have been modified
	// by someone with a pointer to it. If the invariants are no longer upheld,
	// we'll get an error.
	unwrappedPlacesBeen, err := invar.TryUnwrap(unwrappedPerson.PlacesBeen)
	if err != nil {
		panic(err)
	}

	// We know that unwrappedName is non-empty, since its type is NonEmptyString.
	fmt.Println("non-empty name is " + unwrappedName)

	// We know that unwrappedAge is a positive integer, since its type is PositiveInt.
	fmt.Println("positive age is " + fmt.Sprint(unwrappedAge))

	firstPlace := unwrappedPlacesBeen[0]

	unwrappedFirstPlace, err := invar.TryUnwrap(firstPlace)
	if err != nil {
		panic(err)
	}

	fmt.Println("first place is " + unwrappedFirstPlace)
}
```