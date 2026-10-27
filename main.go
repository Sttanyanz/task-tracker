package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"example.com/task-tracker/storage"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return errors.New("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("open PostgreSQL connection: %w", err)
	}
	defer db.Close()

	ctx := context.Background()
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		return fmt.Errorf("connect to PostgreSQL: %w", err)
	}

	st, err := storage.NewPostgresStorage(pingCtx, db)
	if err != nil {
		return err
	}

	t1ID, err := st.CreateTask(ctx, storage.Task{UserID: 133, Name: "Prepare for exams"})
	if err != nil {
		return err
	}
	t2ID, err := st.CreateTask(ctx, storage.Task{UserID: 123, Name: "Clean up the flat"})
	if err != nil {
		return err
	}

	t1, err := st.GetTask(ctx, t1ID)
	if err != nil {
		return err
	}
	fmt.Println(t1)

	t2, err := st.GetTask(ctx, t2ID)
	if err != nil {
		return err
	}
	fmt.Println(t2)

	if err := st.UpdateTask(ctx, t1ID, storage.Task{
		UserID:      133,
		Name:        "Prepare for exams",
		Description: "Math and Physics",
	}); err != nil {
		return err
	}
	t1, err = st.GetTask(ctx, t1ID)
	if err != nil {
		return err
	}
	fmt.Println(t1)

	if err := st.DeleteTask(ctx, t2ID); err != nil {
		return err
	}
	if _, err := st.GetTask(ctx, t2ID); !errors.Is(err, storage.ErrTaskNotFound) {
		if err != nil {
			return err
		}
		return fmt.Errorf("deleted task %d still exists", t2ID)
	}
	fmt.Println(storage.ErrTaskNotFound)
	return nil
}
