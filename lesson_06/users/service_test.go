package users

import (
	"testing"

	"github.com/Nick2603/golang/lesson_06/documentstore"
	"github.com/Nick2603/golang/lesson_06/users/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupService() *Service {
	coll := documentstore.NewCollection(documentstore.CollectionConfig{
		PrimaryKey: "id",
	})
	return NewService(coll)
}

func TestService_CreateUser(t *testing.T) {
	t.Run("creates user successfully", func(t *testing.T) {
		svc := setupService()

		user, err := svc.CreateUser("1", "Alice")
		assert.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, "1", user.ID)
		assert.Equal(t, "Alice", user.Name)
	})

	t.Run("creates multiple users", func(t *testing.T) {
		svc := setupService()

		user1, _ := svc.CreateUser("1", "Alice")
		user2, _ := svc.CreateUser("2", "Bob")

		assert.NotEqual(t, user1.ID, user2.ID)
	})

	t.Run("calls Put on collection store", func(t *testing.T) {
		mockColl := mocks.NewCollectionStore(t)
		svc := NewService(mockColl)

		mockColl.On("Put", mock.AnythingOfType("documentstore.Document")).
			Return(nil).
			Once()

		user, err := svc.CreateUser("1", "Alice")

		assert.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, "1", user.ID)
		assert.Equal(t, "Alice", user.Name)

		mockColl.AssertExpectations(t)
	})
}

func TestService_GetUser(t *testing.T) {
	t.Run("gets existing user", func(t *testing.T) {
		svc := setupService()

		svc.CreateUser("1", "Alice")

		user, err := svc.GetUser("1")
		assert.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, "Alice", user.Name)
	})

	t.Run("returns error for non-existent user", func(t *testing.T) {
		svc := setupService()

		user, err := svc.GetUser("999")
		assert.ErrorIs(t, err, ErrUserNotFound)
		assert.Nil(t, user)
	})

	t.Run("calls Get on collection store", func(t *testing.T) {
		mockColl := mocks.NewCollectionStore(t)
		svc := NewService(mockColl)

		expectedDoc := &documentstore.Document{
			Fields: map[string]documentstore.DocumentField{
				"id": {
					Type:  documentstore.DocumentFieldTypeString,
					Value: "1",
				},
				"name": {
					Type:  documentstore.DocumentFieldTypeString,
					Value: "Alice",
				},
			},
		}

		mockColl.On("Get", "1").
			Return(expectedDoc, nil).
			Once()

		user, err := svc.GetUser("1")

		assert.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, "1", user.ID)
		assert.Equal(t, "Alice", user.Name)

		mockColl.AssertExpectations(t)
	})
}

func TestService_ListUsers(t *testing.T) {
	t.Run("lists all users", func(t *testing.T) {
		svc := setupService()

		svc.CreateUser("1", "Alice")
		svc.CreateUser("2", "Bob")

		users, err := svc.ListUsers()
		assert.NoError(t, err)
		assert.Len(t, users, 2)
	})

	t.Run("returns empty list when no users", func(t *testing.T) {
		svc := setupService()

		users, err := svc.ListUsers()
		assert.NoError(t, err)
		assert.Empty(t, users)
	})

	t.Run("calls List on collection store", func(t *testing.T) {
		mockColl := mocks.NewCollectionStore(t)
		svc := NewService(mockColl)

		docs := []documentstore.Document{
			{
				Fields: map[string]documentstore.DocumentField{
					"id": {
						Type:  documentstore.DocumentFieldTypeString,
						Value: "1",
					},
					"name": {
						Type:  documentstore.DocumentFieldTypeString,
						Value: "Alice",
					},
				},
			},
			{
				Fields: map[string]documentstore.DocumentField{
					"id": {
						Type:  documentstore.DocumentFieldTypeString,
						Value: "2",
					},
					"name": {
						Type:  documentstore.DocumentFieldTypeString,
						Value: "Bob",
					},
				},
			},
		}

		mockColl.On("List").
			Return(docs).
			Once()

		users, err := svc.ListUsers()

		assert.NoError(t, err)
		assert.Len(t, users, 2)

		mockColl.AssertExpectations(t)
	})
}

func TestService_DeleteUser(t *testing.T) {
	t.Run("deletes existing user", func(t *testing.T) {
		svc := setupService()

		svc.CreateUser("1", "Alice")

		err := svc.DeleteUser("1")
		assert.NoError(t, err)

		_, getErr := svc.GetUser("1")
		assert.ErrorIs(t, getErr, ErrUserNotFound)
	})

	t.Run("returns error when deleting non-existent user", func(t *testing.T) {
		svc := setupService()

		err := svc.DeleteUser("999")
		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("calls Delete on collection store", func(t *testing.T) {
		mockColl := mocks.NewCollectionStore(t)
		svc := NewService(mockColl)

		mockColl.On("Delete", "1").
			Return(nil).
			Once()

		err := svc.DeleteUser("1")

		assert.NoError(t, err)

		mockColl.AssertExpectations(t)
	})
}

func TestService_Integration(t *testing.T) {
	t.Run("full CRUD workflow", func(t *testing.T) {
		svc := setupService()

		// Create
		user, err := svc.CreateUser("1", "Alice")
		require.NoError(t, err)
		assert.Equal(t, "Alice", user.Name)

		// Read
		retrieved, err := svc.GetUser("1")
		require.NoError(t, err)
		assert.Equal(t, "Alice", retrieved.Name)

		// List
		users, err := svc.ListUsers()
		require.NoError(t, err)
		assert.Len(t, users, 1)

		// Delete
		err = svc.DeleteUser("1")
		require.NoError(t, err)

		// Verify deletion
		users, err = svc.ListUsers()
		require.NoError(t, err)
		assert.Empty(t, users)
	})
}
