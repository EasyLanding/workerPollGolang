package main

import (
	"fmt"
	"log"
	"sync"
)

func HelloWorld() string {
	fmt.Println("hello World!")
	return `hello world!`
}

type Task struct {
	ID     int
	Data   string
	Status bool
}

type WorkerPool struct {
	Queue        chan *Task
	WorkersCount int
	Result       chan *Task
	wg           sync.WaitGroup
}

func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{
		Queue:        make(chan *Task),
		WorkersCount: workers,
		Result:       make(chan *Task),
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.WorkersCount; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for task := range wp.Queue {
				task.Status = true // Simulating task completion
				wp.Result <- task
			}
		}()
	}
}

func (wp *WorkerPool) AddTask(task *Task) {
	wp.Queue <- task
}

func (wp *WorkerPool) Wait() {
	close(wp.Queue)
	wp.wg.Wait()
	close(wp.Result)
}

func main() {
	pool := NewWorkerPool(5)
	pool.Start()

	tasks := []*Task{
		{ID: 1, Data: "Task 1"},
		{ID: 2, Data: "Task 2"},
		{ID: 3, Data: "Task 3"},
	}

	go func() {
		for _, task := range tasks {
			pool.AddTask(task)
		}
	}()

	go func() {
		for v := range pool.Result {
			log.Printf("Task %d completed: %s status: %v\n", v.ID, v.Data, v.Status)
		}
	}()

	// Waiting for all tasks to complete
	pool.Wait()
}
