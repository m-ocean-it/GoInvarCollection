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
	Name invarcol.NonEmptyString
	Age  invarcol.PositiveInt
}

// ValidPerson is our struct invariant. As an interface, it cannot be directly initialized.
// Also, since the person struct is private, no other package would be able implement that interface.
// The underlying person struct will be accessible via the Invariant.Get method.
type ValidPerson invar.InvariantsHolder[person]

// New is a custom constructor that checks individual field invariants and returns ValidPerson.
// It's also possible to check inter-field invariants within a constructor.
func New(name invarcol.NonEmptyString, age invarcol.PositiveInt) (ValidPerson, error) {
	p := person{Name: name, Age: age}

	// The Inited method, used below, is a way to check whether a certain invariant was initialized.
	// It's important to do, they could be nil, which will cause an error upon accesing the underlying
	// value with TryUnwrap (or panic, if accessing with Unwrap).

	return invar.TryNew(p, []invar.Invariant[person]{
		{
			Name:  "name must be initialized",
			Check: func(p person) bool { return invar.Inited(p.Name) },
		},

		{
			Name:  "age must be initialized",
			Check: func(p person) bool { return invar.Inited(p.Age) },
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

	p, err := person.New(nonEmptyName, positiveAge)
	if err != nil {
		panic(err)
	}

	// It's okay to call Unwrap instead of TryUnwrap here, since we know that the ValidPerson invariant holds up.
	// Otherwise, it would have errored above.
	unwrappedPerson := invar.Unwrap(p)

	// It's okay to call Unwrap instead of TryUnwrap on the structs' fields, since we know that the struct's
	// invariant holds up here and individual field invariants were checked upon instantiation.
	// (Since we just unwrapped the struct and didn't modify it in any way, all it's invariants must hold up.
	// But, to be on the safe side, you can always use TryUnwrap and handle potential errors.)
	unwrappedName := invar.Unwrap(unwrappedPerson.Name)
	unwrappedAge := invar.Unwrap(unwrappedPerson.Age)

	// We know that unwrappedName is non-empty, since it's type is NonEmptyString.
	fmt.Println("non-empty name is " + unwrappedName)

	// We know that unwrappedAge is a positive integer, since it's type is PositiveInt.
	fmt.Println("positive age is " + fmt.Sprint(unwrappedAge))
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