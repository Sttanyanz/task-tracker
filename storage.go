package main

type Storage interface {
	CreateTask(t Task) int
	GetTask(id int) Task
	UpdateTask(id int, t Task)
	DeleteTask(id int)
}
