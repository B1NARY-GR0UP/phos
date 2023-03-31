package main

import (
	"context"
	"fmt"

	"github.com/B1NARY-GR0UP/phos"
)

func hello(_ context.Context, data string) (string, error) {
	return data + "-PHOS", nil
}

func main() {
	ph := phos.New[string]()
	defer close(ph.In)
	ph.AddHandler(hello)
	ph.In <- "BINARY"
	ph.AddHandler(hello)
	res := <-ph.Out
	fmt.Println(res.Data)
}
