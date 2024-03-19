# structTags

structTags allows for marshalling custom and third-party struct tags.

### Installation

Run `go get -u github.com/foresthoffman/structTags`

If you're using `go mod`, run `go mod vendor` afterwards.

### Importing

Import this package by including `github.com/foresthoffman/structTags` in your import block.

e.g.

```go
package main

import(
    ...
    "github.com/foresthoffman/structTags"
)
```

### Usage

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/foresthoffman/structTags"
)

type MyStruct struct {
	Field   string `json:"api_field" custom:"custom_field"`
	Ignored string `json:"ignored" custom:"-"`
}

func main() {
	targetThisTag := "custom"
	ignoreThisTagValue := "-"

	s := MyStruct{
		Field: "some string",
		Ignored: "super-secret-value",
	}

	// custom-tag marshalling...
	m := structTags.NewCustomMarshaller(targetThisTag, ignoreThisTagValue)
	b, err := m.Marshal(s)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
	// Output: {"custom_field":"some string"}

	// versus JSON marshalling...
	b, err = json.Marshal(s)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
	// Output: {"api_field":"some string","ignored":"super-secret-value"}
}
```

### Testing

Run `go test -v -count=1 ./...` in the project root directory. Use the `-count=1` to force the tests to run un-cached.

_That's all, enjoy!_