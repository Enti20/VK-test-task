package workerpool

import (
	"fmt"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	pool := NewWorkerPool(5)
	pool.Start(2)

	// Добавляем задачи
	for i := 1; i <= 5; i++ {
		pool.AddJob(fmt.Sprintf("Тест-%d", i))
	}

	// Добавляем воркеров
	pool.AddWorker()
	pool.AddWorker()

	// Еще задач
	for i := 6; i <= 10; i++ {
		pool.AddJob(fmt.Sprintf("Тест-%d", i))
	}

	time.Sleep(1 * time.Second)
	pool.Stop()
}
