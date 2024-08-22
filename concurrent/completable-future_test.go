package concurrent

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestExecute(t *testing.T) {
	futureA := SupplyAsync(func() (any, error) {
		time.Sleep(2 * time.Second)
		t.Log("future a: done")
		return int64(2233), nil
	})

	futureB := RunAsync(func() error {
		time.Sleep(500 * time.Millisecond)
		t.Log("future b: run async")
		return nil
	})

	holder := ""
	futureC := RunAsync(func() error {
		time.Sleep(1*time.Second + 450*time.Millisecond)
		holder = "bilibili"
		t.Log(`future c: assigned holder by "bilibili"`)
		return nil
	})

	if err := Execute(futureA, futureB, futureC); err != nil {
		t.Fatalf("unexpected error raised: %v", err)
	}

	resultA := futureA.Get()
	if resultA == nil {
		t.Fatalf("future a: result is nil")
	}
	if val, ok := resultA.(int64); !ok || val != int64(2233) {
		t.Fatalf("future a: result is not int64 or is not 2233 in int64")
	}

	if holder != "bilibili" {
		t.Fatalf(`future c: holder is not "bilibili" in string`)
	}
}

func TestRun(t *testing.T) {
	futureA := RunAsync(func() error {
		time.Sleep(2 * time.Second)
		t.Log("future a: done")
		return nil
	})

	futureB := RunAsync(func() error {
		time.Sleep(1 * time.Second)
		t.Log("future b: raise error")
		return errors.New("raise error")
	})

	futureC := SupplyAsync(func() (any, error) {
		time.Sleep(time.Second + 450*time.Millisecond)
		t.Log("future c: return result and raise error")
		return int64(2233), errors.New("bilibili")
	})

	Run(futureA, futureB, futureC)

	if err := futureA.Err(); err != nil {
		t.Fatalf("unexpected error raised from future a: %v", err)
	}

	if err := futureB.Err(); err == nil {
		t.Fatal("expected an error raised from future b, but got nothing")
	}

	resultC, err := futureC.Result()
	if resultC == nil {
		t.Fatalf("expected a result from future c, but got nil")
	}
	if val, ok := resultC.(int64); !ok || val != int64(2233) {
		t.Fatalf("future c: result is not int64 or is not 2233 in int64")
	}
	if err == nil {
		t.Fatalf("expected an error raised from future c, but got nothing")
	}
}

func TestCompletableFuture_Get(t *testing.T) {
	// call Get() multiple times
	future := SupplyAsync(func() (any, error) {
		return 2233, nil
	})
	Run(future)
	numA := future.Get()
	numB := future.Get().(int)
	if numA != numB {
		t.Fatalf("future a: result is not the same as future b")
	}
}

func TestExecute_stabilize_1(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()

			futures := make([]*CompletableFuture, 0)
			for j := 0; j < 100; j++ {
				futures = append(futures, SupplyAsync(func() (any, error) {
					return rand.Intn(100), nil
				}))
			}

			if err := Execute(futures...); err != nil {
				fmt.Println(err)
			}

			for _, future := range futures {
				fmt.Println("iter:", i, ", ret:", future.Get())
			}
		}()
	}
	wg.Wait()
}

func TestExecute_stabilize_2(t *testing.T) {
	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	futures := make([]*CompletableFuture, 0)
	for i := 0; i < 100000; i++ {
		futures = append(futures, SupplyAsync(func() (any, error) {
			rn := rand.Intn(len(numbers))
			return []int{rn, numbers[rn]}, nil
		}))
	}
	if err := Execute(futures...); err != nil {
		t.Fatal(err)
	}

	for _, future := range futures {
		ret := future.Get().([]int)
		if ret[1] != numbers[ret[0]] {
			t.Fatalf("unexpected random number ret: %v", ret)
		}
	}
}
