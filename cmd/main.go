package main

import (
	"fmt"
	"github.com/Enti20/VK-test-task/internal"
	"time"
)

func main() {
	// Создаем пул воркеров с буфером на 10 задач
	pool := internal.NewWorkerPool(10)
	pool.Start(3) // Запускаем пул с 3 воркерами

	// Отправляем первые 10 задач
	for i := 1; i <= 10; i++ {
		zadacha := fmt.Sprintf("Задача-%d", i)
		pool.AddJob(zadacha)
	}

	// Ждем немного, чтобы задачи начали обрабатываться
	time.Sleep(500 * time.Millisecond)

	// Добавляем еще 2 воркера динамически
	fmt.Println("Добавляем еще 2 воркера...")
	pool.AddWorker()
	pool.AddWorker()

	// Отправляем еще 10 задач
	for i := 11; i <= 20; i++ {
		pool.AddJob(fmt.Sprintf("Задача-%d", i))
	}

	// Даем время на обработку всех задач
	time.Sleep(2 * time.Second)

	// Останавливаем пул
	pool.Stop()
	fmt.Println("Все задачи завершены!")
}
