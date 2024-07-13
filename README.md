![PHOS](images/PHOS.png)

PHOS is a channel with internal handlers and diversified options.

[![Go Report Card](https://goreportcard.com/badge/github.com/B1NARY-GR0UP/phos)](https://goreportcard.com/report/github.com/B1NARY-GR0UP/phos)

## Install

```shell
go get github.com/B1NARY-GR0UP/phos
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"

	"github.com/B1NARY-GR0UP/phos"
)

func hello(_ context.Context, data string) (string, error) {
	return data + " by ", nil
}

func world(_ context.Context, data string) (string, error) {
	return data + "BINARY", nil
}

func main() {
	ph := phos.New[string]()
	defer ph.Close()
	ph.Append(hello, world)
	ph.In <- "PHOS"
	res := <-ph.Out
	ph.Delete(ph.Len() - 2)
	fmt.Println(res.Data) // PHOS by BINARY
}
```

## Configuration

| Option               | Default                | Description                                                                          | Example                 |
|----------------------|------------------------|--------------------------------------------------------------------------------------|-------------------------|
| `WithContext`        | `context.Background()` | Set context for PHOS                                                                 | [example](phos_test.go) |
| `WithZero`           | `false`                | Set zero value for return when error happened                                        | [example](phos_test.go) |
| `WithTimeout`        | `3 * time.Second`      | Set timeout for handlers execution                                                   | [example](phos_test.go) |
| `WithErrHandleFunc`  | `nil`                  | Set error handle function for PHOS which will be called when handle error happened   | [example](phos_test.go) |
| `WithErrTimeoutFunc` | `nil`                  | Set error timeout function for PHOS which will be called when timeout error happened | [example](phos_test.go) |
| `WithErrDoneFunc`    | `nil`                  | Set err done function for PHOS which will be called when context done happened       | [example](phos_test.go) |

## Blogs

- [PHOS: A Go channel extension with internal handlers](https://dev.to/justlorain/phos-a-go-channel-extension-with-internal-handlers-4lad) | [中文](https://juejin.cn/post/7216236114981584953)

## License

PHOS is distributed under the [Apache License 2.0](./LICENSE). The licenses of third party dependencies of PHOS are explained [here](./licenses).

## ECOLOGY

<p align="center">
<img src="https://github.com/justlorain/justlorain/blob/main/images/BMS.png" alt="BMS"/>
<br/><br/>
PHOS is a Subproject of the <a href="https://github.com/B1NARY-GR0UP">Basic Middleware Service</a>
</p>