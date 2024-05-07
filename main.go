package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
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

type Order struct {
	ID       int
	Complete bool
}

type sema chan struct{}

var orders []*Order
var completeOrders map[int]bool
var wg sync.WaitGroup
var processTimes chan time.Duration
var sinceProgramStarted time.Duration
var count int
var limitCount int

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

func New(n int) sema {
	return make(sema, n)
}

func (s sema) Inc(k int) {
	for i := 0; i < k; i++ {
		s <- struct{}{}
	}
}

func (s sema) Dec(k int) {
	for i := 0; i < k; i++ {
		<-s
	}
}

func main() {
	// pool := NewWorkerPool(5)
	// pool.Start()

	// tasks := []*Task{
	// 	{ID: 1, Data: "Task 1"},
	// 	{ID: 2, Data: "Task 2"},
	// 	{ID: 3, Data: "Task 3"},
	// }

	// go func() {
	// 	for _, task := range tasks {
	// 		pool.AddTask(task)
	// 	}
	// }()

	// go func() {
	// 	for v := range pool.Result {
	// 		log.Printf("Task %d completed: %s status: %v\n", v.ID, v.Data, v.Status)
	// 	}
	// }()

	// // Waiting for all tasks to complete
	// pool.Wait()

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

	count = 30
	limitCount = 5
	processTimes = make(chan time.Duration, count)
	orders = GenerateOrders(count)
	completeOrders = GenerateCompleteOrders(count)
	programStart := time.Now()
	LimitSpawnOrderProcessing(limitCount)

	wg.Wait()
	sinceProgramStarted = time.Since(programStart)
	go func() {
		time.Sleep(1 * time.Second)
		close(processTimes)
	}()
	checkTimeDifference(limitCount)

	numbers := []int{1, 2, 3, 4, 5}
	n := len(numbers)

	sem := New(n)
	var wg sync.WaitGroup

	for _, num := range numbers {
		wg.Add(1)
		go func(n int) {
			sem.Inc(1)
			defer wg.Done()
			fmt.Println(n)
			sem.Dec(1)
		}(num)
	}

	wg.Wait() // Ждем завершения всех горутин
}

func checkTimeDifference(limitCount int) {
	var averageTime time.Duration
	var orderProcessTotalTime time.Duration
	var orderProcessedCount int
	for v := range processTimes {
		orderProcessedCount++
		orderProcessTotalTime += v
	}
	if orderProcessedCount != count {
		panic("orderProcessedCount != count")
	}
	averageTime = orderProcessTotalTime / time.Duration(orderProcessedCount)
	println("orderProcessTotalTime", orderProcessTotalTime/time.Second)
	println("averageTime", averageTime/time.Second)
	println("sinceProgramStarted", sinceProgramStarted/time.Second)
	println("sinceProgramStarted average", sinceProgramStarted/(time.Duration(orderProcessedCount)*time.Second))
	println("orderProcessTotalTime - sinceProgramStarted", (orderProcessTotalTime-sinceProgramStarted)/time.Second)
	if (orderProcessTotalTime/time.Duration(limitCount)-sinceProgramStarted)/time.Second > 0 {
		panic("(orderProcessTotalTime-sinceProgramStarted)/time.Second > 0")
	}
}

func LimitSpawnOrderProcessing(limitCount int) {
	limit := make(chan struct{}, limitCount)
	for _, order := range orders {
		wg.Add(1)
		limit <- struct{}{}
		go func(order *Order, t time.Time) {
			defer func() { <-limit }()
			OrderProcessing(order, limit, t)
		}(order, time.Now())
	}
}

func OrderProcessing(order *Order, limit chan struct{}, t time.Time) {
	if completeOrders[order.ID] {
		order.Complete = true
	}
	time.Sleep(1 * time.Second)
	select {
	case processTimes <- time.Since(t):
	default:
		// Handle the case when the channel is closed
		fmt.Println("Channel is closed. Cannot send data.")
	}
	wg.Done()
}

func GenerateOrders(count int) []*Order {
	var orders []*Order
	for i := 0; i < count; i++ {
		orders = append(orders, &Order{ID: i})
	}
	return orders
}

func GenerateCompleteOrders(maxOrderID int) map[int]bool {
	completeOrders := make(map[int]bool)
	for i := 0; i < maxOrderID; i++ {
		if rand.Intn(2) == 0 { // 50% chance
			completeOrders[i] = true
		}
	}
	return completeOrders
}
