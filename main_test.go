package main

import (
	"testing"
)

func TestHelloWorld(t *testing.T) {
	expected := "hello world!"
	if result := HelloWorld(); result != expected {
		t.Errorf("HelloWorld() returned %s, expected %s", result, expected)
	}
}

func TestWorkerPool(t *testing.T) {
	pool := NewWorkerPool(2)
	pool.Start()

	task := &Task{ID: 1, Data: "Test Task", Status: false}

	go func() {
		pool.AddTask(task)
	}()

	result := <-pool.Result

	if result.ID != task.ID || result.Data != task.Data || result.Status != true {
		t.Errorf("Task not processed correctly. Expected: %+v, Got: %+v", task, result)
	}

	pool.Wait()
}
