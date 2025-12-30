package documentstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollection_Put(t *testing.T) {
	t.Run("successfully adds document", func(t *testing.T) {
		coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

		doc := Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
				"name": {Type: DocumentFieldTypeString, Value: "Alice"},
			},
		}

		err := coll.Put(doc)
		assert.NoError(t, err)
	})

	t.Run("returns error when fields are nil", func(t *testing.T) {
		coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

		doc := Document{Fields: nil}

		err := coll.Put(doc)
		assert.ErrorIs(t, err, ErrNilValue)
	})

	t.Run("returns error when primary key is missing", func(t *testing.T) {
		coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

		doc := Document{
			Fields: map[string]DocumentField{
				"name": {Type: DocumentFieldTypeString, Value: "Test"},
			},
		}

		err := coll.Put(doc)
		assert.ErrorIs(t, err, ErrInvalidPrimaryKey)
	})

	t.Run("returns error when primary key is not string type", func(t *testing.T) {
		coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

		doc := Document{
			Fields: map[string]DocumentField{
				"id": {Type: DocumentFieldTypeNumber, Value: 123},
			},
		}

		err := coll.Put(doc)
		assert.ErrorIs(t, err, ErrInvalidPrimaryKey)
	})

	t.Run("returns error when primary key value is empty", func(t *testing.T) {
		coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

		doc := Document{
			Fields: map[string]DocumentField{
				"id": {Type: DocumentFieldTypeString, Value: ""},
			},
		}

		err := coll.Put(doc)
		assert.ErrorIs(t, err, ErrInvalidPrimaryKey)
	})
}

func TestCollection_Get(t *testing.T) {
	t.Run("returns document when exists", func(t *testing.T) {
		coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

		doc := Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: "user:1"},
				"name": {Type: DocumentFieldTypeString, Value: "Alice"},
			},
		}

		coll.Put(doc)

		retrieved, err := coll.Get("user:1")
		assert.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, "Alice", retrieved.Fields["name"].Value)
	})

	t.Run("returns error when document not found", func(t *testing.T) {
		coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

		retrieved, err := coll.Get("user:999")
		assert.ErrorIs(t, err, ErrDocumentNotFound)
		assert.Nil(t, retrieved)
	})
}

func TestCollection_Delete(t *testing.T) {
	t.Run("deletes existing document", func(t *testing.T) {
		coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

		doc := Document{
			Fields: map[string]DocumentField{
				"id": {Type: DocumentFieldTypeString, Value: "user:1"},
			},
		}

		coll.Put(doc)

		err := coll.Delete("user:1")
		assert.NoError(t, err)

		_, getErr := coll.Get("user:1")
		assert.ErrorIs(t, getErr, ErrDocumentNotFound)
	})

	t.Run("returns error when document not found", func(t *testing.T) {
		coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

		err := coll.Delete("user:999")
		assert.ErrorIs(t, err, ErrDocumentNotFound)
	})
}

func TestCollection_List(t *testing.T) {
	t.Run("returns all documents", func(t *testing.T) {
		coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

		doc1 := Document{
			Fields: map[string]DocumentField{
				"id": {Type: DocumentFieldTypeString, Value: "user:1"},
			},
		}
		doc2 := Document{
			Fields: map[string]DocumentField{
				"id": {Type: DocumentFieldTypeString, Value: "user:2"},
			},
		}

		coll.Put(doc1)
		coll.Put(doc2)

		docs := coll.List()
		assert.Len(t, docs, 2)
	})

	t.Run("returns empty slice when collection is empty", func(t *testing.T) {
		coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

		docs := coll.List()
		assert.Empty(t, docs)
		assert.NotNil(t, docs)
	})
}
