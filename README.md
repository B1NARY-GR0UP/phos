# PHOS

![]()

> Phosphophyllite

## Install

```shell
go get github.com/B1NARY-GR0UP/phos
```

## Quick Start

### Hello

[example](examples/hello)

```go
package main

import (
	"context"
	"log"
	
	"github.com/B1NARY-GR0UP/phos"
)

func main() {
	ph := phos.New[int](0)
	plusOne := func(ctx context.Context, data int) (int, error) {
		return data + 1, nil
	}
	ph.Handlers = append(ph.Handlers, plusOne)
	ph.In <- 0
	res := <-ph.Out
	log.Printf("res: %d", res.Data)
}
```

## Blogs

- []()

## License

PHOS is distributed under the [Apache License 2.0](./LICENSE). The licenses of third party dependencies of PHOS are explained [here](./licenses).

## END

PHOS is a subproject of the [BINARY WEB ECOLOGY](https://github.com/B1NARY-GR0UP)