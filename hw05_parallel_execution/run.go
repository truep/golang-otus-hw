package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrNullTaskList        = errors.New("there is no tasks")
)

type Task func() error

func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	if len(tasks) == 0 {
		return ErrNullTaskList
	}

	var errCounter int32
	maxErr := int32(m)
	tsChan := make(chan Task)
	wg := &sync.WaitGroup{}

	wg.Add(n)
	for w := 0; w < n; w++ {
		go func() {
			defer wg.Done()
			for t := range tsChan {
				if t != nil {
					if err := t(); err != nil {
						atomic.AddInt32(&errCounter, 1)
					}
				}
			}
		}()
	}

	for _, t := range tasks {
		if atomic.LoadInt32(&errCounter) == maxErr {
			break
		}
		tsChan <- t
	}

	close(tsChan)
	wg.Wait()

	if errCounter > maxErr {
		return ErrErrorsLimitExceeded
	}

	return nil
}
