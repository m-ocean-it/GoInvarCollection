# Collection of [GoInvar](https://github.com/m-ocean-it/GoInvar) types

## Example usage

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
}
```