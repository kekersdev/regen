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
    xeger := regen.FromPattern(`0x[\da-f]{16}`)
    fmt.Println(xeger.MustGenerate())
}
```
`FromPattern()` accepts a regular expression as a single parameter, parses it assuming it uses Perl-like syntax and returns a ready-to-use instance of string generator or `nil` in case of an error.

`MustGenerate()` method returns a new generated string. Will panic if an error occurs during generation.

### Advanced
```golang
package main

import (
    "fmt"
    "regexp/syntax"

    "github.com/kekersdev/regen"
)

func main() {
    xeger := &XegerGenerator

    xeger.AddPattern(`0x[\da-f]{16}`,     syntax.Perl,  false)
    xeger.AddPattern(`#[[:alpha:]]{5,6}`, syntax.POSIX, true)

    xeger.SetUnboundLimit(10)

    if str, err := xeger.Generate(); err != nil {
        fmt.Println("Failed to generate string")
    } else {
        fmt.Println(str)
    }
}
```
`AddPattern()` method offers more flexibility as it allows to specify expression parser flags and wether or not parsed expression should be simplified. _Note: when a generator instance holds multiple expressions, a random one will be selected when generating a string._

`SetUnboundLimit()` method allows to change the default maximum number of repetitions for patterns with unbound quantifiers (`+`/`*`). _Note: default limit for repetitions is `32`, the same as for the original **regen**._

`Generate()` method works the same as `MustGenerate()` but returns an error alongside the generated string. If an error occurs during generation the returned string will be empty.
_Note: while the original **regen** may panic in certain situations, `Generate()` will return an error instead._

License
-------
regen is licensed under a 2-clause BSD license. This can be found in LICENSE.txt.
