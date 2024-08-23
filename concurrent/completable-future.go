package concurrent

import "sync/atomic"

func RunAsync(f func() error) *CompletableFuture[any] {
	future := &CompletableFuture[any]{
		r: &runnable{run: f},
	}
	go future.execute()
	return future
}

// Wait all runnable *CompletableFuture[T] done, futures must defined in same Type
func Wait[T any](futures ...*CompletableFuture[T]) {
	if len(futures) == 0 {
		return
	}
	for i := range futures {
		futures[i].Wait()
	}
}

func SupplyAsync[T any](f func() (T, error)) *CompletableFuture[T] {
	future := &CompletableFuture[T]{
		s:          &supplier[T]{get: f},
		resultChan: make(chan T, 1),
	}
	go future.execute()
	return future
}

type CompletableFuture[T any] struct {
	s          *supplier[T] // supplier future
	resultChan chan T       // result channel
	result     *T           // result holder

	r *runnable // runnable future

	done atomic.Bool
	err  error // error holder
}

func (future *CompletableFuture[T]) Result() (T, error) {
	result := future.Get()
	err := future.Err()
	return result, err
}

func (future *CompletableFuture[T]) Get() T {
	if future.result != nil {
		return *future.result
	}

	future.Wait()

	var empty T
	if len(future.resultChan) == 0 {
		return empty
	}

	result, ok := <-future.resultChan
	if !ok {
		return empty
	}
	future.result = &result

	future.close() // close the channels immediately after reading and caching the result

	return result
}

func (future *CompletableFuture[T]) Err() error {
	future.Wait()
	return future.err
}

func (future *CompletableFuture[T]) IsDone() bool {
	return future.done.Load()
}

func (future *CompletableFuture[T]) Wait() {
	for {
		if future.done.Load() {
			break
		}
	}
}

func (future *CompletableFuture[T]) execute() {
	defer future.done.Store(true)

	if future == nil || (future.r == nil && future.s == nil) {
		future.close()
		return
	}

	s := future.s
	r := future.r

	if s != nil {
		// supplier
		resultChan := future.resultChan
		res, err := s.get()
		resultChan <- res
		if err != nil {
			future.err = err
			return
		}
	}

	if r != nil {
		// runnable
		if err := r.run(); err != nil {
			future.err = err
			return
		}
	}
}

func (future *CompletableFuture[T]) close() {
	if future == nil {
		return
	}
	if future.resultChan != nil {
		close(future.resultChan)
	}
}

type supplier[T any] struct {
	get func() (T, error)
}

type runnable struct {
	run func() error
}
