# Collection of [GoInvar](https://github.com/m-ocean-it/GoInvar) types

## Implemented types

- [Positive integer](https://github.com/m-ocean-it/GoInvarCollection/blob/main/positive_int.go)
- [Non-empty string](https://github.com/m-ocean-it/GoInvarCollection/blob/main/non_empty_string.go)
- [Non-empty slice](https://github.com/m-ocean-it/GoInvarCollection/blob/main/non_empty_slice.go)

### TODO
- [ ] non-empty map
- [ ] positive int64
- [ ] positive int32
- [ ] positive float64
- [ ] positive float32
- [ ] non-nil error
- [ ] non-nil pointer
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

#### Approach 1: use `Invariant[T]`, where `T` is the struct

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

// ValidPerson is our invariants holder. As an interface, it cannot be directly initialized.
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
	p := person{Name: name, Age: age, PlacesBeen: placesBeen}

	// The Inited method, used below, is a way to check whether a certain invariant holder was initialized.
	// It's important to do, since they could be nil (if they were passed as nil to the constructor or you
	// forgot to set some field when initializing the person struct). If you fail to check for that, it
	// will cause an error upon accesing the underlying value with TryUnwrap (or panic, if accessing with Unwrap).

	return invar.TryNew(p, []invar.Invariant[person]{
		{
			Name:  "name must be initialized",
			Check: func(p person) bool { return invar.Inited(p.Name) },
		},
		{
			Name:  "age must be initialized",
			Check: func(p person) bool { return invar.Inited(p.Age) },
		},
		{
			Name:  "placesBeen must be initialized",
			Check: func(p person) bool { return invar.Inited(p.PlacesBeen) },
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

	// Notice the level of nesting! ValidPerson will hold a NonEmptySlice which will hold a NonEmptyString.
	// Each element is an InvariantsHolder which satisfies the invariants across its lifetime.
	placesBeen := invarcol.NewNonEmptySlice([]invarcol.NonEmptyString{
		invarcol.NewNonEmptyString("London"),
	})

	p, err := person.New(nonEmptyName, positiveAge, placesBeen)
	if err != nil {
		panic(err)
	}

	// It's okay to call Unwrap instead of TryUnwrap here, since we know that the ValidPerson invariant holds up.
	// Otherwise, it would have errored above.
	unwrappedPerson := invar.Unwrap(p)

	// It's okay to call Unwrap instead of TryUnwrap on the structs' fields, since we know that the struct's
	// invariant holds up here. (Since we just unwrapped the struct and didn't modify it in any way.
	// But, to be on the safe side, you can always use TryUnwrap and handle potential errors.)
	unwrappedName := invar.Unwrap(unwrappedPerson.Name)
	unwrappedAge := invar.Unwrap(unwrappedPerson.Age)
	unwrappedPlacesBeen := invar.Unwrap(unwrappedPerson.PlacesBeen)

	// We know that unwrappedName is non-empty, since its type is NonEmptyString.
	fmt.Println("non-empty name is " + unwrappedName)

	// We know that unwrappedAge is a positive integer, since its type is PositiveInt.
	fmt.Println("positive age is " + fmt.Sprint(unwrappedAge))

	// Accessing the first place is safe since we know that unwrappedPlacesBeen is non-empty,
	// because its type is NonEmptySlice. We, also, need to unwrap the NonEmptyString we're getting.
	fmt.Println("first place is " + invar.Unwrap(unwrappedPlacesBeen[0]))
}
```

#### Approach 2: create accessor-methods and unwrap fields within them

`person/person.go`:
```go
package person

import (
    invar "github.com/m-ocean-it/GoInvar"
    invarcol "github.com/m-ocean-it/GoInvarCollection"
)

type person struct {
    name invarcol.NonEmptyString
    age  invarcol.PositiveInt
}

type ValidPerson interface {
    GetName() string
    GetAge() int
	
    self() *person // to disallow other implementors
}

func New(name string, age int) (ValidPerson, error) {
    nonEmptyName, err := invarcol.TryNewNonEmptyString(name)
    if err != nil {
        return nil, errors.New("")
    }

    positiveAge, err := invarcol.TryNewPositiveInt(age)
    if err != nil {
	    return nil, errors.New("")
    }

    return &person{
        name: nonEmptyName,
        age:  positiveAge,
    }, nil
}

func (p *person) GetName() string {
    return invar.Unwrap(p.name)
}

func (p *person) GetAge() int {
    return invar.Unwrap(p.age)
}

func (p *person) self() *person {
    return p
}
```
`main.go`:
```go
import (
    "app/person"
    "fmt"
    invar "github.com/m-ocean-it/GoInvar"
)

func main() {
	p, err := person.New("Simon", 29)
	if err != nil {
		panic(err)
	}

	fmt.Printf("name is %s\n", p.GetName())
	fmt.Printf("age is %d\n", p.GetAge())
}
```