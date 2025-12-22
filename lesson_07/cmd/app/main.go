package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Nick2603/golang/lesson_07/internal/documentstore"
	"github.com/Nick2603/golang/lesson_07/internal/users"
)

func main() {
	// Configure structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	fmt.Println("=== Document Store with Dump/Restore ===")

	// Create store
	store := documentstore.NewStoreWithLogger(logger)

	usersColl, _ := store.CreateCollection(
		"users",
		&documentstore.CollectionConfig{PrimaryKey: "id"},
	)

	userService := users.NewService(usersColl)

	fmt.Println("Creating users...")
	userService.CreateUser("1", "Alice")
	userService.CreateUser("2", "Bob")
	userService.CreateUser("3", "Charlie")

	fmt.Println("\nListing users...")
	list, _ := userService.ListUsers()
	fmt.Println(list)

	// Dump to file
	filename := "store_backup.json"
	fmt.Printf("\nDumping store to file: %s\n", filename)
	if err := store.DumpToFile(filename); err != nil {
		fmt.Printf("Error dumping: %v\n", err)
		return
	}

	fmt.Println("\nDeleting user 2...")
	userService.DeleteUser("2")

	fmt.Println("\nUsers after deletion:")
	list, _ = userService.ListUsers()
	fmt.Println(list)

	// Restore from file
	fmt.Printf("\nRestoring store from file: %s\n", filename)
	restoredStore, err := documentstore.NewStoreFromFile(filename)
	if err != nil {
		fmt.Printf("Error restoring: %v\n", err)
		return
	}

	restoredUsersColl, _ := restoredStore.GetCollection("users")
	restoredUserService := users.NewService(restoredUsersColl)

	fmt.Println("\nUsers after restore (should have all 3 users again):")
	list, _ = restoredUserService.ListUsers()
	fmt.Println(list)

	// Cleanup
	os.Remove(filename)
	fmt.Println("\n=== Demo completed ===")
}
