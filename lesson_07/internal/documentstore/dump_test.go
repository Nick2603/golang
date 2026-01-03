package documentstore

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_Dump(t *testing.T) {
	t.Run("dumps empty store", func(t *testing.T) {
		store := NewStore()

		data, err := store.Dump()
		assert.NoError(t, err)
		assert.NotEmpty(t, data)
	})

	t.Run("dumps store with collections and documents", func(t *testing.T) {
		store := NewStore()

		coll, _ := store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
		doc := Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
				"name": {Type: DocumentFieldTypeString, Value: "Alice"},
			},
		}
		coll.Put(doc)

		data, err := store.Dump()
		assert.NoError(t, err)
		assert.NotEmpty(t, data)
	})
}

func TestNewStoreFromDump(t *testing.T) {
	t.Run("restores store from dump", func(t *testing.T) {
		// Create original store
		originalStore := NewStore()
		coll, _ := originalStore.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})

		doc1 := Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
				"name": {Type: DocumentFieldTypeString, Value: "Alice"},
			},
		}
		doc2 := Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: "user:2"},
				"name": {Type: DocumentFieldTypeString, Value: "Bob"},
			},
		}
		coll.Put(doc1)
		coll.Put(doc2)

		// Dump
		data, err := originalStore.Dump()
		require.NoError(t, err)

		// Restore
		restoredStore, err := NewStoreFromDump(data)
		require.NoError(t, err)

		// Verify
		restoredColl, err := restoredStore.GetCollection("users")
		require.NoError(t, err)

		docs := restoredColl.List()
		assert.Len(t, docs, 2)

		retrievedDoc, err := restoredColl.Get("user:1")
		require.NoError(t, err)
		assert.Equal(t, "Alice", retrievedDoc.Fields["name"].Value)
	})

	t.Run("returns error for invalid dump", func(t *testing.T) {
		invalidData := []byte("invalid json")

		store, err := NewStoreFromDump(invalidData)
		assert.Error(t, err)
		assert.Nil(t, store)
	})
}

func TestStore_DumpToFile(t *testing.T) {
	t.Run("dumps store to file", func(t *testing.T) {
		store := NewStore()
		coll, _ := store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})

		doc := Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
				"name": {Type: DocumentFieldTypeString, Value: "Alice"},
			},
		}
		coll.Put(doc)

		filename := "test_dump.json"
		defer os.Remove(filename)

		err := store.DumpToFile(filename)
		assert.NoError(t, err)

		// Verify file exists
		_, err = os.Stat(filename)
		assert.NoError(t, err)
	})
}

func TestNewStoreFromFile(t *testing.T) {
	t.Run("loads store from file", func(t *testing.T) {
		// Create and dump original store
		originalStore := NewStore()
		coll, _ := originalStore.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})

		doc := Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
				"name": {Type: DocumentFieldTypeString, Value: "Alice"},
				"age":  {Type: DocumentFieldTypeNumber, Value: int64(25)},
			},
		}
		coll.Put(doc)

		filename := "test_restore.json"
		defer os.Remove(filename)

		err := originalStore.DumpToFile(filename)
		require.NoError(t, err)

		// Load from file
		restoredStore, err := NewStoreFromFile(filename)
		require.NoError(t, err)

		// Verify
		restoredColl, err := restoredStore.GetCollection("users")
		require.NoError(t, err)

		retrievedDoc, err := restoredColl.Get("user:1")
		require.NoError(t, err)
		assert.Equal(t, "Alice", retrievedDoc.Fields["name"].Value)
		assert.Equal(t, float64(25), retrievedDoc.Fields["age"].Value)
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		store, err := NewStoreFromFile("nonexistent.json")
		assert.Error(t, err)
		assert.Nil(t, store)
	})
}

func TestDumpAndRestore_Integration(t *testing.T) {
	t.Run("full dump and restore workflow", func(t *testing.T) {
		// Create original store with multiple collections
		originalStore := NewStore()

		usersColl, _ := originalStore.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
		productsColl, _ := originalStore.CreateCollection("products", &CollectionConfig{PrimaryKey: "sku"})

		// Add users
		usersColl.Put(Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
				"name": {Type: DocumentFieldTypeString, Value: "Alice"},
			},
		})
		usersColl.Put(Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: "user:2"},
				"name": {Type: DocumentFieldTypeString, Value: "Bob"},
			},
		})

		// Add products
		productsColl.Put(Document{
			Fields: map[string]DocumentField{
				"sku":   {Type: DocumentFieldTypeString, Value: "prod:1"},
				"title": {Type: DocumentFieldTypeString, Value: "Laptop"},
				"price": {Type: DocumentFieldTypeNumber, Value: int64(1000)},
			},
		})

		// Dump to file
		filename := "test_integration.json"
		defer os.Remove(filename)

		err := originalStore.DumpToFile(filename)
		require.NoError(t, err)

		// Restore from file
		restoredStore, err := NewStoreFromFile(filename)
		require.NoError(t, err)

		// Verify users collection
		restoredUsersColl, err := restoredStore.GetCollection("users")
		require.NoError(t, err)
		userDocs := restoredUsersColl.List()
		assert.Len(t, userDocs, 2)

		// Verify products collection
		restoredProductsColl, err := restoredStore.GetCollection("products")
		require.NoError(t, err)
		productDocs := restoredProductsColl.List()
		assert.Len(t, productDocs, 1)

		// Verify specific document
		laptop, err := restoredProductsColl.Get("prod:1")
		require.NoError(t, err)
		assert.Equal(t, "Laptop", laptop.Fields["title"].Value)
		assert.Equal(t, float64(1000), laptop.Fields["price"].Value)
	})
}
