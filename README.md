# PHOS

[![Go Report Card](https://goreportcard.com/badge/github.com/B1NARY-GR0UP/phos)](https://goreportcard.com/report/github.com/B1NARY-GR0UP/phos) [![coveralls](https://coveralls.io/repos/B1NARY-GR0UP/phos/badge.svg?branch=main&service=github)](https://coveralls.io/github/B1NARY-GR0UP/phos?branch=main)

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
	"fmt"

	"github.com/B1NARY-GR0UP/phos"
)

func hello(_ context.Context, data string) (string, error) {
	return data + " PHOS", nil
}

func main() {
	ph := phos.New[string]()
	ph.Handlers = append(ph.Handlers, hello)
	ph.In <- "BINARY"
	res := <-ph.Out
	fmt.Println(res.Data)
}
```

## Blogs

- []()

## License

PHOS is distributed under the [Apache License 2.0](./LICENSE). The licenses of third party dependencies of PHOS are explained [here](./licenses).

## END

PHOS is a subproject of the [BINARY WEB ECOLOGY](https://github.com/B1NARY-GR0UP)