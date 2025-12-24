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

		coll, err := store.CreateCollection("users", &CollectionConfig{
			PrimaryKey: "id",
		})

		assert.NoError(t, err)
		require.NotNil(t, coll)
		assert.Equal(t, "id", coll.cfg.PrimaryKey)
	})

	t.Run("returns error when collection already exists", func(t *testing.T) {
		store := NewStore()

		store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})

		coll, err := store.CreateCollection("users", &CollectionConfig{
			PrimaryKey: "key",
		})

		assert.ErrorIs(t, err, ErrCollectionAlreadyExists)
		assert.Nil(t, coll)
	})

	t.Run("returns error when config is nil", func(t *testing.T) {
		store := NewStore()

		coll, err := store.CreateCollection("users", nil)

		assert.ErrorIs(t, err, ErrNilValue)
		assert.Nil(t, coll)
	})
}

func TestStore_GetCollection(t *testing.T) {
	t.Run("retrieves existing collection", func(t *testing.T) {
		store := NewStore()

		store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})

		coll, err := store.GetCollection("users")
		assert.NoError(t, err)
		require.NotNil(t, coll)
		assert.Equal(t, "id", coll.cfg.PrimaryKey)
	})

	t.Run("returns error for non-existent collection", func(t *testing.T) {
		store := NewStore()

		coll, err := store.GetCollection("users")
		assert.ErrorIs(t, err, ErrCollectionNotFound)
		assert.Nil(t, coll)
	})
}

func TestStore_DeleteCollection(t *testing.T) {
	t.Run("deletes existing collection", func(t *testing.T) {
		store := NewStore()

		store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})

		err := store.DeleteCollection("users")
		assert.NoError(t, err)

		_, getErr := store.GetCollection("users")
		assert.ErrorIs(t, getErr, ErrCollectionNotFound)
	})

	t.Run("returns error when deleting non-existent collection", func(t *testing.T) {
		store := NewStore()

		err := store.DeleteCollection("users")
		assert.ErrorIs(t, err, ErrCollectionNotFound)
	})
}
