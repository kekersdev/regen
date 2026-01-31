regen (package)
=======
This is a fork of [nilium/regen](https://github.com/nilium/regen) transformed into a golang package.

# Usage
```console
    $ go get github.com/kekersdev/regen
```
## Examples
### Simple
```golang
package main

import (
    "fmt"

    "github.com/kekersdev/regen"
)

func main() {
    xeger := regen.NewGenerator(`0x[\da-f]{16}`)
    fmt.Println(xeger.MustGenerate())
}
```
`NewGenerator()` accepts a regular expression as a single parameter, parses it assuming it uses Perl-like syntax and returns a ready-to-use instance of string generator or `nil` in case of an error.

`MustGenerate()` method returns a new generated string. Will panic if something goes wrong during generation.

### Advanced
```golang
package main

import (
    "fmt"
    "regexp/syntax"

    "github.com/kekersdev/regen"
)

func main() {
    xeger, err := regen.NewGeneratorAdvanced(
        []string{
            `[[:alpha:]]\.[[:digit:]]{1,3}`,    // expression
            `[[:alpha:]]{2,8}`,                 // expression
        },
        syntax.POSIX,                           // parser flags
        true,                                   // option to simplify expression
    )

    if err != nil {
        fmt.Println("Failed to initialize string generator")
        return
    }

    // limit for unbound repetitions can be customized
    xeger.SetUnboundLimit(10)

    if str, err := xeger.Generate(); err != nil {
        fmt.Println("Failed to generate string")
    } else {
        fmt.Println(str)
    }
}
```
`NewGeneratorAdvanced()` offers more flexibility as it allows to specify multiple expressions simultaneously, expression parser flags and wether or not the parsed expressions should be simplified. _Note: when multiple expressions are provided, a random one will be selected for each new string generation._

`SetUnboundLimit()` method allows to change the default maximum number of repetitions for patterns with unbound quantifiers (`+`/`*`). _Note: default limit for repetitions is `32`, the same as in the original `regen`._

`Generate()` method works the same as `MustGenerate()` but returns an error alongside the generated string. If an error occurs during generation the returned string will be empty.
_Note: while the original may panic in certain situations, `Generate()` will instead return an error in such cases._

License
-------
regen is licensed under a 2-clause BSD license. This can be found in LICENSE.txt.
