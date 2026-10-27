package storage

type Task struct {
	ID          int64
	UserID      int64
	Name        string
	Description string
	Completed   bool
}
