package users

import (
	"testing"

	"github.com/Nick2603/golang/lesson_07/internal/documentstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupService(t *testing.T) *Service {
	t.Helper() // Mark this as a test helper

	coll := documentstore.NewCollection(documentstore.CollectionConfig{
		PrimaryKey: "id",
	})

	return NewService(coll)
}

func TestService_CreateUser(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		userName  string
		wantErr   error
		checkUser func(*testing.T, *User)
	}{
		{
			name:     "creates user successfully",
			id:       "1",
			userName: "Alice",
			wantErr:  nil,
			checkUser: func(t *testing.T, user *User) {
				assert.Equal(t, "1", user.ID)
				assert.Equal(t, "Alice", user.Name)
			},
		},
		{
			name:     "creates user with different name",
			id:       "2",
			userName: "Bob",
			wantErr:  nil,
			checkUser: func(t *testing.T, user *User) {
				assert.Equal(t, "2", user.ID)
				assert.Equal(t, "Bob", user.Name)
			},
		},
		{
			name:     "creates user with empty name",
			id:       "3",
			userName: "",
			wantErr:  nil,
			checkUser: func(t *testing.T, user *User) {
				assert.Equal(t, "3", user.ID)
				assert.Equal(t, "", user.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupService(t)

			user, err := svc.CreateUser(tt.id, tt.userName)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, user)
				if tt.checkUser != nil {
					tt.checkUser(t, user)
				}
			}
		})
	}
}

func TestService_GetUser(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(*Service)
		userID    string
		wantErr   error
		wantName  string
	}{
		{
			name: "gets existing user",
			setupFunc: func(svc *Service) {
				svc.CreateUser("1", "Alice")
			},
			userID:   "1",
			wantErr:  nil,
			wantName: "Alice",
		},
		{
			name: "gets user from multiple users",
			setupFunc: func(svc *Service) {
				svc.CreateUser("1", "Alice")
				svc.CreateUser("2", "Bob")
				svc.CreateUser("3", "Charlie")
			},
			userID:   "2",
			wantErr:  nil,
			wantName: "Bob",
		},
		{
			name:      "returns error for non-existent user",
			setupFunc: func(svc *Service) {},
			userID:    "999",
			wantErr:   ErrUserNotFound,
			wantName:  "",
		},
		{
			name: "returns error after user is deleted",
			setupFunc: func(svc *Service) {
				svc.CreateUser("1", "Alice")
				svc.DeleteUser("1")
			},
			userID:   "1",
			wantErr:  ErrUserNotFound,
			wantName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupService(t)
			tt.setupFunc(svc)

			user, err := svc.GetUser(tt.userID)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, tt.wantName, user.Name)
			}
		})
	}
}

func TestService_ListUsers(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(*Service)
		wantCount int
		checkFunc func(*testing.T, []User)
	}{
		{
			name: "lists all users",
			setupFunc: func(svc *Service) {
				svc.CreateUser("1", "Alice")
				svc.CreateUser("2", "Bob")
			},
			wantCount: 2,
			checkFunc: nil,
		},
		{
			name:      "returns empty list when no users",
			setupFunc: func(svc *Service) {},
			wantCount: 0,
			checkFunc: func(t *testing.T, users []User) {
				assert.Empty(t, users)
				assert.NotNil(t, users)
			},
		},
		{
			name: "lists users after some are deleted",
			setupFunc: func(svc *Service) {
				svc.CreateUser("1", "Alice")
				svc.CreateUser("2", "Bob")
				svc.CreateUser("3", "Charlie")
				svc.DeleteUser("2")
			},
			wantCount: 2,
			checkFunc: nil,
		},
		{
			name: "lists single user",
			setupFunc: func(svc *Service) {
				svc.CreateUser("1", "Alice")
			},
			wantCount: 1,
			checkFunc: func(t *testing.T, users []User) {
				assert.Equal(t, "Alice", users[0].Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupService(t)
			tt.setupFunc(svc)

			users, err := svc.ListUsers()

			assert.NoError(t, err)
			assert.Len(t, users, tt.wantCount)
			if tt.checkFunc != nil {
				tt.checkFunc(t, users)
			}
		})
	}
}

func TestService_DeleteUser(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(*Service)
		userID    string
		wantErr   error
	}{
		{
			name: "deletes existing user",
			setupFunc: func(svc *Service) {
				svc.CreateUser("1", "Alice")
			},
			userID:  "1",
			wantErr: nil,
		},
		{
			name:      "returns error when deleting non-existent user",
			setupFunc: func(svc *Service) {},
			userID:    "999",
			wantErr:   ErrUserNotFound,
		},
		{
			name: "returns error when deleting already deleted user",
			setupFunc: func(svc *Service) {
				svc.CreateUser("1", "Alice")
				svc.DeleteUser("1")
			},
			userID:  "1",
			wantErr: ErrUserNotFound,
		},
		{
			name: "deletes one user from multiple users",
			setupFunc: func(svc *Service) {
				svc.CreateUser("1", "Alice")
				svc.CreateUser("2", "Bob")
				svc.CreateUser("3", "Charlie")
			},
			userID:  "2",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupService(t)
			tt.setupFunc(svc)

			err := svc.DeleteUser(tt.userID)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)

				// Verify deletion
				_, getErr := svc.GetUser(tt.userID)
				assert.ErrorIs(t, getErr, ErrUserNotFound)
			}
		})
	}
}

// Keep integration test as is - it's already clear
func TestService_Integration(t *testing.T) {
	t.Run("full CRUD workflow", func(t *testing.T) {
		svc := setupService(t)

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

func TestService_ConcurrentOperations(t *testing.T) {
	t.Run("creates multiple users without conflicts", func(t *testing.T) {
		svc := setupService(t)

		// Create multiple users
		users := []struct {
			id   string
			name string
		}{
			{"1", "Alice"},
			{"2", "Bob"},
			{"3", "Charlie"},
			{"4", "Diana"},
			{"5", "Eve"},
		}

		for _, u := range users {
			_, err := svc.CreateUser(u.id, u.name)
			require.NoError(t, err)
		}

		// Verify all users exist
		list, err := svc.ListUsers()
		require.NoError(t, err)
		assert.Len(t, list, len(users))

		// Verify each user can be retrieved
		for _, u := range users {
			retrieved, err := svc.GetUser(u.id)
			require.NoError(t, err)
			assert.Equal(t, u.name, retrieved.Name)
		}
	})
}
