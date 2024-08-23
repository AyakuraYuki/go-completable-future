# go-concurrent

CompletableFuture in Golang.

A simple implementation of CompletableFuture's basic `RunAsync` and `SupplyAsync` features using `sync.WaitGroup`, `goroutine`, `select` and `channel`.

## Usage

```go
package main

import (
	"github.com/AyakuraYuki/go-concurrent/concurrent"
	"github.com/AyakuraYuki/go-concurrent/futuretask"
)

func main() {
	// do something
}

```

### futuretask

> For more detail, please see [`example/example-futuretask.go`](https://github.com/AyakuraYuki/go-concurrent/blob/main/example/example-futuretask.go)

```go
package demo

import (
	"fmt"
	"log"

	"github.com/AyakuraYuki/go-concurrent/futuretask"
)

func demoExecute() {
	runAsync := futuretask.PlanRun(func() error {
		// do something
		return nil
	})

	supplyAsync := futuretask.PlanSupply(func() (any, error) {
		// do something
		return 0, nil
	})

	if err := futuretask.Execute(runAsync, supplyAsync); err != nil {
		log.Fatal(err)
	}

	result := supplyAsync.Get().(int)
	fmt.Println("result:", result)
}

func demoRun() {
	futureA := futuretask.PlanSupply(func() (any, error) {
		// do something
		return "bilibili", nil
	})

	holder := ""
	futureB := futuretask.PlanRun(func() error {
		holder = "2233"
		return nil
	})

	futuretask.Run(futureA, futureB)

	resultA, err := futureA.Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("result a:", resultA)
	fmt.Println("holder:", holder)
}

```

### concurrent

> For more detail, please see [`example/example-completable-future.go`](https://github.com/AyakuraYuki/go-concurrent/blob/main/example/example-completable-future.go)

```go
package demo

import (
	"errors"
	"fmt"
	"github.com/AyakuraYuki/go-concurrent/concurrent"
	"log"
)

func demoRunAsync() {
	holder := ""
	taskA := concurrent.RunAsync(func() error {
		// do something
		holder = "2233"
		return nil
	})
	taskB := concurrent.RunAsync(func() error {
		// raise an error
		return errors.New("error")
	})

	concurrent.Wait(taskA, taskB)

	fmt.Println(holder)
	if err := taskB.Err(); err != nil {
		log.Fatal(err)
	}
}

type Foo struct {
	Val string `json:"val"`
}

func demoSupplyAsync() {
	taskA := concurrent.SupplyAsync(func() (int64, error) {
		return int64(2233), nil
	})
	taskB := concurrent.SupplyAsync(func() (string, error) {
		return "bilibili", errors.New("error")
	})

	taskC := concurrent.SupplyAsync(func() (*Foo, error) {
		return &Foo{Val: "66"}, nil
	})

	number := taskA.Get()
	fmt.Printf("number: %d\n", number)

	err := taskB.Err()
	fmt.Printf("err: %v\n", err)

	result, err := taskC.Result()
	fmt.Printf("result of Foo: %#v\n", result)
}

```
