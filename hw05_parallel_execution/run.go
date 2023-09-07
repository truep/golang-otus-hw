package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var errCounter int32
	tsChan := make(chan Task)
	wg := &sync.WaitGroup{}

	wg.Add(n)
	for w := 0; w < n; w++ {
		go func() {
			defer wg.Done()
			for t := range tsChan {
				if err := t(); err != nil {
					atomic.AddInt32(&errCounter, 1)
				}
			}
		}()
	}

	for _, t := range tasks {
		if atomic.LoadInt32(&errCounter) >= int32(m) {
			break
		}
		tsChan <- t
	}

	close(tsChan)
	wg.Wait()

	if errCounter > int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
