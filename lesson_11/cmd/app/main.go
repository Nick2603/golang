package main

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/Nick2603/golang/lesson_11/internal/documentstore"
	"github.com/Nick2603/golang/lesson_11/internal/users"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)

	store := documentstore.NewStoreWithLogger(logger)

	usersColl, _ := store.CreateCollection("users", &documentstore.CollectionConfig{PrimaryKey: "id"})

	userService := users.NewService(usersColl)

	var wg sync.WaitGroup

	for i := range 1000 {

		wg.Go(func() {
			userID := fmt.Sprintf("%d", i)

			userService.CreateUser(userID, fmt.Sprintf("User%d", i))

			userService.GetUser(userID)

			userService.ListUsers()

			userService.DeleteUser(userID)
		})
	}

	wg.Wait()

	fmt.Println("All goroutines completed")
}
