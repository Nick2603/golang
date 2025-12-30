package main

import (
	"fmt"

	"github.com/Nick2603/golang/lesson_05/documentstore"
	"github.com/Nick2603/golang/lesson_05/users"
)

func main() {
	store := documentstore.NewStore()

	usersColl, _ := store.CreateCollection(
		"users",
		&documentstore.CollectionConfig{PrimaryKey: "id"},
	)

	userService := users.NewService(usersColl)

	fmt.Println("Creating users...")
	userService.CreateUser("1", "Alice")
	userService.CreateUser("2", "Bob")

	fmt.Println("\nListing users...")
	list, _ := userService.ListUsers()
	fmt.Println(list)

	fmt.Println("\nGetting user 1...")
	u, _ := userService.GetUser("1")
	fmt.Println(u)

	fmt.Println("\nDeleting user 1...")
	userService.DeleteUser("1")

	fmt.Println("\nFinal users list:")
	list, _ = userService.ListUsers()
	fmt.Println(list)
}
