package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type ErrorsLimit struct {
	errorsCount    int
	maxErrorsCount int
	mutext         sync.RWMutex
}

func (el *ErrorsLimit) IncrementErrorsCount() {
	el.mutext.Lock()
	defer el.mutext.Unlock()
	el.errorsCount++
}

func (el *ErrorsLimit) LimitExceeded() bool {
	el.mutext.RLock()
	defer el.mutext.RUnlock()
	return el.errorsCount >= el.maxErrorsCount
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workersCount, maxErrorsCount int) error {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	taskChannel := make(chan Task)
	defer close(taskChannel)

	el := ErrorsLimit{maxErrorsCount: maxErrorsCount}

	wg.Add(workersCount)
	consume := func() {
		defer wg.Done()
		for task := range taskChannel {
			err := task()
			if err != nil {
				el.IncrementErrorsCount()
			}
		}
	}

	for i := 0; i < workersCount; i++ {
		go consume()
	}

	for _, task := range tasks {
		if el.LimitExceeded() {
			return ErrErrorsLimitExceeded
		}
		taskChannel <- task
	}
	return nil
}
