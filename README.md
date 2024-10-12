# Collection of [GoInvar](https://github.com/m-ocean-it/GoInvar) types

## Example usage

```go
package main

import (
	"fmt"

	goinvar "github.com/m-ocean-it/GoInvar"
	goinvarcol "github.com/m-ocean-it/GoInvarCollection"
)

func main() {
	positiveInt, err := goinvarcol.TryNewPositiveInt(2)
	if err != nil {
		panic(err)
	}

	fmt.Println("positive int is " + fmt.Sprint(goinvar.Unwrap(positiveInt)))
}
```