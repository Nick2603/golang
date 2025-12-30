package documentstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollection_Put(t *testing.T) {
	t.Run("successfully adds document with valid primary key", func(t *testing.T) {
		collection := NewCollection(CollectionConfig{PrimaryKey: "id"})

		doc := Document{
			Fields: map[string]DocumentField{
				"id": {
					Type:  DocumentFieldTypeString,
					Value: "user:1",
				},
				"name": {
					Type:  DocumentFieldTypeString,
					Value: "Іван Петренко",
				},
			},
		}

		err := collection.Put(doc)
		assert.NoError(t, err)

		retrieved, found := collection.Get("user:1")
		require.True(t, found)
		assert.Equal(t, "Іван Петренко", retrieved.Fields["name"].Value)
	})

	t.Run("returns error when fields are nil", func(t *testing.T) {
		collection := NewCollection(CollectionConfig{PrimaryKey: "id"})

		doc := Document{Fields: nil}

		err := collection.Put(doc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "fields cannot be nil")
	})

	t.Run("returns error when primary key is missing", func(t *testing.T) {
		collection := NewCollection(CollectionConfig{PrimaryKey: "id"})

		doc := Document{
			Fields: map[string]DocumentField{
				"name": {
					Type:  DocumentFieldTypeString,
					Value: "Test",
				},
			},
		}

		err := collection.Put(doc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must contain primary key field 'id'")
	})

	t.Run("returns error when primary key is not string type", func(t *testing.T) {
		collection := NewCollection(CollectionConfig{PrimaryKey: "id"})

		doc := Document{
			Fields: map[string]DocumentField{
				"id": {
					Type:  DocumentFieldTypeNumber,
					Value: 123,
				},
			},
		}

		err := collection.Put(doc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be of type string")
	})

	t.Run("returns error when primary key value is empty string", func(t *testing.T) {
		collection := NewCollection(CollectionConfig{PrimaryKey: "id"})

		doc := Document{
			Fields: map[string]DocumentField{
				"id": {
					Type:  DocumentFieldTypeString,
					Value: "",
				},
			},
		}

		err := collection.Put(doc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "non-empty string")
	})

	t.Run("updates existing document", func(t *testing.T) {
		collection := NewCollection(CollectionConfig{PrimaryKey: "id"})

		doc1 := Document{
			Fields: map[string]DocumentField{
				"id": {
					Type:  DocumentFieldTypeString,
					Value: "user:1",
				},
				"age": {
					Type:  DocumentFieldTypeNumber,
					Value: 25,
				},
			},
		}

		doc2 := Document{
			Fields: map[string]DocumentField{
				"id": {
					Type:  DocumentFieldTypeString,
					Value: "user:1",
				},
				"age": {
					Type:  DocumentFieldTypeNumber,
					Value: 30,
				},
			},
		}

		collection.Put(doc1)
		collection.Put(doc2)

		retrieved, found := collection.Get("user:1")
		require.True(t, found)
		assert.Equal(t, 30, retrieved.Fields["age"].Value)
	})
}

func TestCollection_Get(t *testing.T) {
	t.Run("returns document when key exists", func(t *testing.T) {
		collection := NewCollection(CollectionConfig{PrimaryKey: "key"})

		doc := Document{
			Fields: map[string]DocumentField{
				"key": {
					Type:  DocumentFieldTypeString,
					Value: "user:1",
				},
				"name": {
					Type:  DocumentFieldTypeString,
					Value: "Марія",
				},
			},
		}

		collection.Put(doc)

		retrieved, found := collection.Get("user:1")
		assert.True(t, found)
		require.NotNil(t, retrieved)
		assert.Equal(t, "Марія", retrieved.Fields["name"].Value)
	})

	t.Run("returns false when key does not exist", func(t *testing.T) {
		collection := NewCollection(CollectionConfig{PrimaryKey: "key"})

		retrieved, found := collection.Get("user:999")
		assert.False(t, found)
		assert.Nil(t, retrieved)
	})
}

func TestCollection_Delete(t *testing.T) {
	t.Run("deletes existing document", func(t *testing.T) {
		collection := NewCollection(CollectionConfig{PrimaryKey: "key"})

		doc := Document{
			Fields: map[string]DocumentField{
				"key": {
					Type:  DocumentFieldTypeString,
					Value: "user:1",
				},
			},
		}

		collection.Put(doc)

		deleted := collection.Delete("user:1")
		assert.True(t, deleted)

		_, found := collection.Get("user:1")
		assert.False(t, found)
	})

	t.Run("returns false when deleting non-existent document", func(t *testing.T) {
		collection := NewCollection(CollectionConfig{PrimaryKey: "key"})

		deleted := collection.Delete("user:999")
		assert.False(t, deleted)
	})
}

func TestCollection_List(t *testing.T) {
	t.Run("returns all documents", func(t *testing.T) {
		collection := NewCollection(CollectionConfig{PrimaryKey: "key"})

		doc1 := Document{
			Fields: map[string]DocumentField{
				"key":  {Type: DocumentFieldTypeString, Value: "user:1"},
				"name": {Type: DocumentFieldTypeString, Value: "Іван"},
			},
		}

		doc2 := Document{
			Fields: map[string]DocumentField{
				"key":  {Type: DocumentFieldTypeString, Value: "user:2"},
				"name": {Type: DocumentFieldTypeString, Value: "Марія"},
			},
		}

		collection.Put(doc1)
		collection.Put(doc2)

		docs := collection.List()
		assert.Len(t, docs, 2)
	})

	t.Run("returns empty slice when collection is empty", func(t *testing.T) {
		collection := NewCollection(CollectionConfig{PrimaryKey: "key"})

		docs := collection.List()
		assert.Empty(t, docs)
		assert.NotNil(t, docs)
	})
}
