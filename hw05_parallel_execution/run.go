package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	tasksCh := make(chan Task)
	errorsCh := make(chan error)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(ctx, tasksCh, errorsCh, &wg)
	}

	go func() {
		defer close(tasksCh)

		for _, task := range tasks {
			if ctx.Err() != nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			case tasksCh <- task:
			}
		}
	}()

	isError := make(chan bool, 1)
	go func() {
		defer close(isError)

		errorsCount := 0
		for range errorsCh {
			errorsCount++
			if errorsCount == m {
				isError <- true
				cancel()
			}
		}
	}()

	wg.Wait()
	close(errorsCh)

	if <-isError {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(ctx context.Context, tasks <-chan Task, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-tasks:
			if !ok {
				return
			}
			err := task()
			if err != nil {
				errors <- err
			}
		}
	}
}
