package futuretask

import "sync"

func PlanRun(f func() error) *Task {
	return &Task{
		r: &runnable{run: f},
	}
}

func PlanSupply(f func() (any, error)) *Task {
	return &Task{
		s:          &supplier{get: f},
		resultChan: make(chan any, 1),
	}
}

// Execute given Task, returns error if one of the Task resulting error.
func Execute(futures ...*Task) error {
	if len(futures) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(futures))

	doneChan := make(chan bool)
	errChan := make(chan error)

	for _, future := range futures {

		go func(future *Task) {
			defer wg.Done()

			s := future.s
			r := future.r

			if s != nil {
				// supplier
				result := future.resultChan
				res, err := s.get()
				result <- res
				if err != nil {
					future.err = err
					errChan <- err
					return
				}
			}

			if r != nil {
				// runnable
				if err := r.run(); err != nil {
					future.err = err
					errChan <- err
					return
				}
			}

		}(future)

	}

	go func() {
		wg.Wait()
		close(doneChan)
		close(errChan)
	}()

	select {
	case err := <-errChan:
		return err
	case <-doneChan:
		return nil
	}
}

// Run given Task without returning error, you should handle error manually in each Task.
func Run(futures ...*Task) {
	if len(futures) == 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(futures))

	doneChan := make(chan bool)

	for _, future := range futures {

		go func(future *Task) {
			defer wg.Done()

			s := future.s
			r := future.r

			if s != nil {
				// supplier
				result := future.resultChan
				res, err := s.get()
				result <- res
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

		}(future)

	}

	go func() {
		wg.Wait()
		close(doneChan)
	}()

	select {
	case <-doneChan:
		return
	}
}

type Task struct {
	s          *supplier // supplier future
	resultChan chan any  // result channel
	result     any       // result holder

	r *runnable // runnable future

	err error // error holder
}

func (future *Task) Result() (any, error) {
	result := future.Get()
	err := future.Err()
	return result, err
}

func (future *Task) Get() any {
	if future.result != nil {
		return future.result
	}
	var empty any
	if len(future.resultChan) == 0 {
		return empty
	}
	result, ok := <-future.resultChan
	if !ok {
		return empty
	}
	future.result = result
	close(future.resultChan)
	return future.result
}

func (future *Task) Err() error {
	return future.err
}

type supplier struct {
	get func() (any, error)
}

type runnable struct {
	run func() error
}
