package storage

import "context"

type Storage interface {
	CreateTask(ctx context.Context, task Task) (int64, error)
	GetTask(ctx context.Context, id int64) (Task, error)
	UpdateTask(ctx context.Context, id int64, task Task) error
	DeleteTask(ctx context.Context, id int64) error
}
