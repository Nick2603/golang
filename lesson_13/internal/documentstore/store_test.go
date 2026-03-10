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
	tests := []struct {
		name      string
		collName  string
		cfg       *CollectionConfig
		setupFunc func(*Store) // For testing duplicate names
		wantErr   error
	}{
		{
			name:     "creates new collection successfully",
			collName: "users",
			cfg:      &CollectionConfig{PrimaryKey: "id"},
			wantErr:  nil,
		},
		{
			name:     "returns error when collection already exists",
			collName: "users",
			cfg:      &CollectionConfig{PrimaryKey: "key"},
			setupFunc: func(s *Store) {
				s.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
			},
			wantErr: ErrCollectionAlreadyExists,
		},
		{
			name:     "returns error when config is nil",
			collName: "users",
			cfg:      nil,
			wantErr:  ErrNilValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewStore()

			// Setup
			if tt.setupFunc != nil {
				tt.setupFunc(store)
			}

			// Test
			coll, err := store.CreateCollection(tt.collName, tt.cfg)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, coll)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, coll)
				assert.Equal(t, tt.cfg.PrimaryKey, coll.cfg.PrimaryKey)
			}
		})
	}
}

func TestStore_GetCollection(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(*Store)
		collName  string
		wantErr   error
		wantPK    string
	}{
		{
			name: "retrieves existing collection",
			setupFunc: func(s *Store) {
				s.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
			},
			collName: "users",
			wantErr:  nil,
			wantPK:   "id",
		},
		{
			name:      "returns error for non-existent collection",
			setupFunc: func(s *Store) {},
			collName:  "users",
			wantErr:   ErrCollectionNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewStore()
			tt.setupFunc(store)

			coll, err := store.GetCollection(tt.collName)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, coll)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, coll)
				assert.Equal(t, tt.wantPK, coll.cfg.PrimaryKey)
			}
		})
	}
}

func TestStore_DeleteCollection(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(*Store)
		collName  string
		wantErr   error
	}{
		{
			name: "deletes existing collection",
			setupFunc: func(s *Store) {
				s.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
			},
			collName: "users",
			wantErr:  nil,
		},
		{
			name:      "returns error when deleting non-existent collection",
			setupFunc: func(s *Store) {},
			collName:  "users",
			wantErr:   ErrCollectionNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewStore()
			tt.setupFunc(store)

			err := store.DeleteCollection(tt.collName)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)

				// Verify deletion
				_, getErr := store.GetCollection(tt.collName)
				assert.ErrorIs(t, getErr, ErrCollectionNotFound)
			}
		})
	}
}
