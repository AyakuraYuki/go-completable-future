# go-completable-future

CompletableFuture in Golang.

A simple implementation of CompletableFuture's basic `RunAsync` and `SupplyAsync` features using `sync.WaitGroup`, `goroutine`, `select` and `channel`.

## Usage

> For more detail, please see [`example/example.go`](https://github.com/AyakuraYuki/go-completable-future/blob/main/example/example.go)

```go
package demo

import (
	"fmt"
	"log"

	"github.com/AyakuraYuki/go-completable-future/concurrent"
)

func demoExecute() {
	runAsync := concurrent.RunAsync(func() error {
		// do something
		return nil
	})

	supplyAsync := concurrent.SupplyAsync(func() (any, error) {
		// do something
		return 0, nil
	})

	if err := concurrent.Execute(runAsync, supplyAsync); err != nil {
		log.Fatal(err)
	}

	result := supplyAsync.Get().(int)
	fmt.Println("result:", result)
}

func demoRun() {
	futureA := concurrent.SupplyAsync(func() (any, error) {
		// do something
		return "bilibili", nil
	})

	holder := ""
	futureB := concurrent.RunAsync(func() error {
		holder = "2233"
		return nil
	})

	concurrent.Run(futureA, futureB)

	resultA, err := futureA.Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("result a:", resultA)
	fmt.Println("holder:", holder)
}

```
