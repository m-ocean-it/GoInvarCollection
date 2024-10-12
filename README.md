# Collection of [GoInvar](https://github.com/m-ocean-it/GoInvar) types

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

    fmt.Println("positive int is " + fmt.Sprint(invar.Unwrap(positiveInt)))
    // or use TryUnwrap
}
```

### Encoding struct invariants

#### Approach 1: wrap whole struct, then unwrap it, then unwrap separate fields

`person/person.go`:
```go
package person

import (
    invar "github.com/m-ocean-it/GoInvar"
    invarcol "github.com/m-ocean-it/GoInvarCollection"
)

type person struct {
    Name invarcol.NonEmptyString
    Age  invarcol.PositiveInt
}

type ValidPerson invar.Invariant[person]

func New(name string, age int) (ValidPerson, error) {
    nonEmptyName, err := invarcol.TryNewNonEmptyString(name)
    if err != nil {
        return nil, errors.New("name is invalid")
    }
    positiveAge, err := invarcol.TryNewPositiveInt(age)
    if err != nil {
        return nil, errors.New("age is invalid")
    }

    p := person{Name: nonEmptyName, Age: positiveAge}

    return invar.TryNew(p, []invar.Condition[person]{
        func(p person) bool { return invar.Inited(p.Name) },
        func(p person) bool { return invar.Inited(p.Age) },
    })
}
```

`main.go`:
```go
import (
    "app/person"
    "fmt"
    invar "github.com/m-ocean-it/GoInvar"
)

wrappedPerson, err := person.New("Simon", 29)
if err != nil {
    panic(err)
}

validPerson := invar.Unwrap(wrappedPerson) // or use TryUnwrap

fmt.Printf("name is %s\n", invar.Unwrap(validPerson.Name))
fmt.Printf("age is %d\n", invar.Unwrap(validPerson.Age))
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