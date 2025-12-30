package documentstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStore(t *testing.T) {
	store := NewStore()

	assert.NotNil(t, store)
	assert.NotNil(t, store.collections)
}

func TestStore_CreateCollection(t *testing.T) {
	t.Run("creates new collection successfully", func(t *testing.T) {
		store := NewStore()

		created, collection := store.CreateCollection("users", &CollectionConfig{
			PrimaryKey: "id",
		})

		assert.True(t, created)
		require.NotNil(t, collection)
		assert.Equal(t, "id", collection.cfg.PrimaryKey)
	})

	t.Run("returns false when collection already exists", func(t *testing.T) {
		store := NewStore()

		store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})

		created, collection := store.CreateCollection("users", &CollectionConfig{
			PrimaryKey: "key",
		})

		assert.False(t, created)
		assert.Nil(t, collection)
	})

	t.Run("returns false when config is nil", func(t *testing.T) {
		store := NewStore()

		created, collection := store.CreateCollection("users", nil)

		assert.False(t, created)
		assert.Nil(t, collection)
	})

	t.Run("creates multiple different collections", func(t *testing.T) {
		store := NewStore()

		created1, col1 := store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
		created2, col2 := store.CreateCollection("products", &CollectionConfig{PrimaryKey: "sku"})

		assert.True(t, created1)
		assert.True(t, created2)
		assert.NotNil(t, col1)
		assert.NotNil(t, col2)
	})
}

func TestStore_GetCollection(t *testing.T) {
	t.Run("retrieves existing collection", func(t *testing.T) {
		store := NewStore()

		store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})

		collection, exists := store.GetCollection("users")
		assert.True(t, exists)
		require.NotNil(t, collection)
		assert.Equal(t, "id", collection.cfg.PrimaryKey)
	})

	t.Run("returns false for non-existent collection", func(t *testing.T) {
		store := NewStore()

		collection, exists := store.GetCollection("users")
		assert.False(t, exists)
		assert.Nil(t, collection)
	})
}

func TestStore_DeleteCollection(t *testing.T) {
	t.Run("deletes existing collection", func(t *testing.T) {
		store := NewStore()

		store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})

		deleted := store.DeleteCollection("users")
		assert.True(t, deleted)

		_, exists := store.GetCollection("users")
		assert.False(t, exists)
	})

	t.Run("returns false when deleting non-existent collection", func(t *testing.T) {
		store := NewStore()

		deleted := store.DeleteCollection("users")
		assert.False(t, deleted)
	})
}

func TestStore_Integration(t *testing.T) {
	t.Run("full workflow: create, add documents, query, delete", func(t *testing.T) {
		store := NewStore()

		// Create collection
		created, users := store.CreateCollection("users", &CollectionConfig{
			PrimaryKey: "id",
		})
		require.True(t, created)
		require.NotNil(t, users)

		// Add documents
		doc1 := Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
				"name": {Type: DocumentFieldTypeString, Value: "Іван"},
				"age":  {Type: DocumentFieldTypeNumber, Value: 25},
			},
		}

		doc2 := Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: "user:2"},
				"name": {Type: DocumentFieldTypeString, Value: "Марія"},
				"age":  {Type: DocumentFieldTypeNumber, Value: 30},
			},
		}

		err1 := users.Put(doc1)
		err2 := users.Put(doc2)
		assert.NoError(t, err1)
		assert.NoError(t, err2)

		// List documents
		docs := users.List()
		assert.Len(t, docs, 2)

		// Get specific document
		retrieved, found := users.Get("user:1")
		require.True(t, found)
		assert.Equal(t, "Іван", retrieved.Fields["name"].Value)

		// Delete document
		deleted := users.Delete("user:2")
		assert.True(t, deleted)

		// Verify deletion
		docs = users.List()
		assert.Len(t, docs, 1)
	})
}
