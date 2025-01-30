package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error // Определяем тип Task как функцию, которая возвращает ошибку.

// Run запускает задачи в n горутинах и прекращает работу при получении m ошибок от задач.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	tasksCh := make(chan Task)   // Канал для передачи задач в горутины.
	errorsCh := make(chan error) // Канал для передачи ошибок от горутин.

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(ctx, tasksCh, errorsCh, &wg)
	}

	// Запускаем горутину для передачи задач в tasksCh.
	go func() {
		defer close(tasksCh) // Закрываем tasksCh по завершении функции.

		// Проходим по всем задачам и отправляем их в tasksCh.
		for _, task := range tasks {
			if ctx.Err() != nil { // Проверяем, не отменен ли контекст.
				return // Если отменен, выходим из функции.
			}
			select {
			case <-ctx.Done(): // Проверяем, не отменен ли контекст.
				return // Если отменен, выходим из функции.
			case tasksCh <- task: // Отправляем задачу в tasksCh.
			}
		}
	}()

	isError := make(chan bool, 1) // Канал для информирования о превышении лимита ошибок.
	go func() {
		defer close(isError) // Закрываем isError по завершении функции.

		errorsCount := 0     // Счетчик ошибок.
		for range errorsCh { // Проходим по всем ошибкам из errorsCh.
			errorsCount++         // Увеличиваем счетчик ошибок.
			if errorsCount == m { // Если достигнут лимит ошибок m,
				isError <- true // Отправляем сигнал в isError.
				cancel()        // Отменяем контекст, чтобы остановить работу горутин.
			}
		}
	}()

	wg.Wait()
	close(errorsCh) // Закрываем errorsCh.

	if <-isError { // Проверяем, был ли сигнал из канала о превышении лимита ошибок.
		return ErrErrorsLimitExceeded // Возвращаем ошибку о превышении лимита ошибок.
	}

	return nil // Если ошибок не было или лимит не превышен, возвращаем nil.
}

// worker выполняет задачи из tasksCh и отправляет ошибки в errorsCh.
func worker(ctx context.Context, tasks <-chan Task, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done() // Уменьшаем счетчик WaitGroup при завершении функции.

	for {
		select {
		case <-ctx.Done(): // Проверяем, не отменен ли контекст.
			return // Если отменен, выходим из функции.
		case task, ok := <-tasks: // Получаем задачу из tasksCh.
			if !ok { // Если tasksCh закрыт,
				return // выходим из функции.
			}
			err := task()   // Выполняем задачу.
			if err != nil { // Если произошла ошибка,
				errors <- err // Отправляем ошибку в errorsCh.
			}
		}
	}
}
