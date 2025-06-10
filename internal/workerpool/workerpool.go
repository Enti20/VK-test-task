package workerpool

import (
	"fmt"
	"sync"
	"time"
)

// Worker представляет одного воркера
type Worker struct {
	ID    int
	tasks chan string // Канал для получения задач
	quit  chan struct{}
	wg    *sync.WaitGroup
}

// NewWorker создает нового воркера
func NewWorker(id int, wg *sync.WaitGroup) *Worker {
	return &Worker{
		ID:    id,
		tasks: make(chan string),
		quit:  make(chan struct{}),
		wg:    wg,
	}
}

// Start запускает воркера
func (w *Worker) Start() {
	go func() {
		defer w.wg.Done()
		fmt.Printf("Воркер %d начал работу!\n", w.ID)
		for {
			select {
			case zadacha, ok := <-w.tasks:
				if !ok {
					fmt.Printf("Воркер %d остановлен, канал закрыт\n", w.ID)
					return
				}
				fmt.Printf("Воркер %d обрабатывает: %s\n", w.ID, zadacha)
				time.Sleep(50 * time.Millisecond) // Имитация работы
			case <-w.quit:
				fmt.Printf("Воркер %d получил сигнал остановки\n", w.ID)
				return
			}
		}
	}()
}

// Stop останавливает воркера
func (w *Worker) Stop() {
	close(w.quit)
}

// WorkerPool управляет пулом воркеров
type WorkerPool struct {
	kolvoZadach chan string     // Канал для задач
	workers     map[int]*Worker // Список воркеров
	workerID    int             // Счетчик ID воркеров
	mu          sync.Mutex      // Мьютекс для безопасности
	quit        chan struct{}   // Сигнал остановки
	wg          sync.WaitGroup  // Группа ожидания
	workerChans []chan string   // Список каналов воркеров
}

// NewWorkerPool создает новый пул
func NewWorkerPool(bufferSize int) *WorkerPool {
	return &WorkerPool{
		kolvoZadach: make(chan string, bufferSize),
		workers:     make(map[int]*Worker),
		quit:        make(chan struct{}),
		workerChans: make([]chan string, 0),
	}
}

// Start запускает пул с заданным количеством воркеров
func (wp *WorkerPool) Start(kolvo int) {
	for i := 0; i < kolvo; i++ {
		wp.addWorker()
	}
	go wp.distributeTasks()
}

// addWorker добавляет нового воркера
func (wp *WorkerPool) addWorker() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	wp.workerID++
	newWorker := NewWorker(wp.workerID, &wp.wg)
	wp.workers[newWorker.ID] = newWorker
	wp.workerChans = append(wp.workerChans, newWorker.tasks)
	wp.wg.Add(1)
	newWorker.Start()
	fmt.Printf("Добавлен воркер %d, всего: %d\n", newWorker.ID, len(wp.workers))
}

// AddWorker добавляет еще одного воркера динамически
func (wp *WorkerPool) AddWorker() {
	wp.addWorker()
}

// AddJob добавляет задачу в пул
func (wp *WorkerPool) AddJob(zadacha string) {
	select {
	case wp.kolvoZadach <- zadacha:
	case <-wp.quit:
		fmt.Printf("Пул остановлен, задача %s не добавлена\n", zadacha)
	}
}

// distributeTasks распределяет задачи между воркерами
func (wp *WorkerPool) distributeTasks() {
	index := 0
	for {
		select {
		case zadacha, ok := <-wp.kolvoZadach:
			if !ok {
				fmt.Println("Канал задач закрыт, распределение остановлено")
				return
			}
			wp.mu.Lock()
			if len(wp.workerChans) == 0 {
				fmt.Printf("Нет воркеров для задачи: %s\n", zadacha)
				wp.mu.Unlock()
				continue
			}
			targetChan := wp.workerChans[index]
			index = (index + 1) % len(wp.workerChans)
			wp.mu.Unlock()
			select {
			case targetChan <- zadacha:
			case <-wp.quit:
				fmt.Printf("Пул остановлен, задача %s не распределена\n", zadacha)
				return
			}
		case <-wp.quit:
			fmt.Println("Пул получил сигнал остановки")
			return
		}
	}
}

// Stop останавливает пул
func (wp *WorkerPool) Stop() {
	close(wp.quit)
	close(wp.kolvoZadach)
	wp.wg.Wait()
	wp.mu.Lock()
	defer wp.mu.Unlock()
	for _, w := range wp.workers {
		close(w.tasks)
	}
	fmt.Println("Пул воркеров остановлен")
}
