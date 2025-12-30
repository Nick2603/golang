package main

import (
	"fmt"

	"github.com/Nick2603/golang/lesson_04/documentstore"
)

func main() {
	fmt.Println("=== Document Store Test Scenario ===")

	store := documentstore.NewStore()

	created, users := store.CreateCollection(
		"users",
		&documentstore.CollectionConfig{
			PrimaryKey: "key",
		},
	)

	if !created {
		panic("failed to create users collection")
	}

	fmt.Println("✓ Collection 'users' created")

	fmt.Println("1. Adding documents...")

	user1 := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"key": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "user:1",
			},
			"name": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "Іван Петренко",
			},
			"age": {
				Type:  documentstore.DocumentFieldTypeNumber,
				Value: 25,
			},
			"active": {
				Type:  documentstore.DocumentFieldTypeBool,
				Value: true,
			},
		},
	}

	user2 := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"key": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "user:2",
			},
			"name": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "Марія Іваненко",
			},
			"age": {
				Type:  documentstore.DocumentFieldTypeNumber,
				Value: 30,
			},
			"active": {
				Type:  documentstore.DocumentFieldTypeBool,
				Value: false,
			},
			"tags": {
				Type:  documentstore.DocumentFieldTypeArray,
				Value: []string{"admin", "moderator"},
			},
		},
	}

	product := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"key": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "product:100",
			},
			"title": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "Ноутбук",
			},
			"price": {
				Type:  documentstore.DocumentFieldTypeNumber,
				Value: 25000.50,
			},
			"inStock": {
				Type:  documentstore.DocumentFieldTypeBool,
				Value: true,
			},
		},
	}

	if err := users.Put(user1); err != nil {
		fmt.Printf("Error adding user1: %v\n", err)
	} else {
		fmt.Println("✓ Added user:1")
	}

	if err := users.Put(user2); err != nil {
		fmt.Printf("Error adding user2: %v\n", err)
	} else {
		fmt.Println("✓ Added user:2")
	}

	if err := users.Put(product); err != nil {
		fmt.Printf("Error adding product: %v\n", err)
	} else {
		fmt.Println("✓ Added product:100")
	}

	fmt.Println("\n2. Getting documents...")
	if doc, found := users.Get("user:1"); found {
		fmt.Printf(
			"✓ Found user:1 - Name: %v, Age: %v\n",
			doc.Fields["name"].Value,
			doc.Fields["age"].Value,
		)
	} else {
		fmt.Println("✗ user:1 not found")
	}

	if _, found := users.Get("user:999"); !found {
		fmt.Println("✓ user:999 not found (as expected)")
	}

	fmt.Println("\n3. Listing all documents...")
	allDocs := users.List()
	fmt.Printf("Total documents: %d\n", len(allDocs))
	for _, doc := range allDocs {
		fmt.Printf("  - %v\n", doc.Fields["key"].Value)
	}

	fmt.Println("\n4. Deleting documents...")
	if users.Delete("user:2") {
		fmt.Println("✓ Deleted user:2")
	} else {
		fmt.Println("✗ Failed to delete user:2")
	}

	if !users.Delete("user:999") {
		fmt.Println("✓ user:999 not deleted (doesn't exist)")
	}

	fmt.Println("\n5. Listing after deletion...")
	allDocs = users.List()
	fmt.Printf("Total documents: %d\n", len(allDocs))
	for _, doc := range allDocs {
		fmt.Printf("  - %v\n", doc.Fields["key"].Value)
	}

	fmt.Println("\n6. Testing validation...")

	invalidDoc1 := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"name": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "Test",
			},
		},
	}

	if err := users.Put(invalidDoc1); err != nil {
		fmt.Printf("✓ Validation works: %v\n", err)
	}

	invalidDoc2 := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"key": {
				Type:  documentstore.DocumentFieldTypeNumber,
				Value: 123,
			},
		},
	}

	if err := users.Put(invalidDoc2); err != nil {
		fmt.Printf("✓ Type validation works: %v\n", err)
	}

	fmt.Println("\n=== Test completed ===")
}
