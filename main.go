package main

import (
	"fmt"
	"log"
	"sync"
)

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

type User struct {
	ID      int
	Name    string
	Bonuses int
}

type BonusOperation struct {
	UserID int
	Amount int
}

var users = []*User{
	{ID: 1, Name: "Bob", Bonuses: 100},
	{ID: 2, Name: "Alice", Bonuses: 200},
	{ID: 3, Name: "Kate", Bonuses: 300},
	{ID: 4, Name: "Tom", Bonuses: 400},
	{ID: 5, Name: "John", Bonuses: 500},
}

// DeductBonuses - вычитает бонусы у пользователя
func DeductBonuses(users []*User, bonusesOperations []BonusOperation) {
	done := make(chan bool)

	for i := range users {
		go func(index int) {
			// здесь мы как будто обращаемся во внешний сервис для списания бонусов
			users[index].Bonuses -= bonusesOperations[index].Amount
			done <- true
		}(i)
	}

	// Ждем завершения всех горутин
	for range users {
		<-done
	}
}

func HelloWorld() string {
	fmt.Println("hello World!")
	return `hello world!`
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

	bonusOperations := []BonusOperation{
		{UserID: 1, Amount: 10},
		{UserID: 2, Amount: 20},
		{UserID: 3, Amount: 30},
		{UserID: 4, Amount: 40},
		{UserID: 5, Amount: 50},
	}

	DeductBonuses(users, bonusOperations)

	for _, user := range users {
		// fmt.Printf
		fmt.Printf("User %s has %d bonuses\n", user.Name, user.Bonuses)
	}
}
