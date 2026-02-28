package services

import (
	"context"
	"testing"
	"time"

	"github.com/Nick2603/golang/final_project/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func setupUserService(mt *mtest.T) *UserService {
	mt.Helper()
	return NewUserService(mt.Coll)
}

func TestUserService_Create(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	tests := []struct {
		name         string
		input        CreateUserInput
		addResponses func(mt *mtest.T)
		wantErr      error
		checkUser    func(*testing.T, *models.User)
	}{
		{
			name: "creates user successfully",
			input: CreateUserInput{
				Username: "alice",
				Email:    "alice@example.com",
				Password: "password123",
			},
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(
					mtest.CreateCursorResponse(0, ns, mtest.FirstBatch),
					mtest.CreateSuccessResponse(),
				)
			},
			wantErr: nil,
			checkUser: func(t *testing.T, u *models.User) {
				assert.Equal(t, "alice", u.Username)
				assert.Equal(t, "alice@example.com", u.Email)
				assert.NotEmpty(t, u.Password)
				assert.NotEqual(t, "password123", u.Password)
				assert.False(t, u.ID.IsZero())
			},
		},
		{
			name: "creates user with different credentials",
			input: CreateUserInput{
				Username: "bob",
				Email:    "bob@example.com",
				Password: "bobspassword",
			},
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(
					mtest.CreateCursorResponse(0, ns, mtest.FirstBatch),
					mtest.CreateSuccessResponse(),
				)
			},
			wantErr: nil,
			checkUser: func(t *testing.T, u *models.User) {
				assert.Equal(t, "bob", u.Username)
				assert.Equal(t, "bob@example.com", u.Email)
			},
		},
		{
			name: "returns error when email already exists",
			input: CreateUserInput{
				Username: "alice",
				Email:    "alice@example.com",
				Password: "password123",
			},
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(
					mtest.CreateCursorResponse(0, ns, mtest.FirstBatch, bson.D{{Key: "n", Value: int32(1)}}),
				)
			},
			wantErr: ErrEmailAlreadyExists,
		},
	}

	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			tt.addResponses(mt)

			svc := setupUserService(mt)
			user, err := svc.Create(context.Background(), tt.input)

			if tt.wantErr != nil {
				assert.ErrorIs(mt.T, err, tt.wantErr)
				assert.Nil(mt.T, user)
			} else {
				assert.NoError(mt.T, err)
				require.NotNil(mt.T, user)
				if tt.checkUser != nil {
					tt.checkUser(mt.T, user)
				}
			}
		})
	}
}

func TestUserService_FindByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	userID := primitive.NewObjectID()
	now := time.Now()

	userDoc := bson.D{
		{Key: "_id", Value: userID},
		{Key: "username", Value: "alice"},
		{Key: "email", Value: "alice@example.com"},
		{Key: "password", Value: "$2a$10$hashedpassword"},
		{Key: "createdAt", Value: now},
		{Key: "updatedAt", Value: now},
	}

	tests := []struct {
		name         string
		id           primitive.ObjectID
		addResponses func(mt *mtest.T)
		wantErr      error
		checkUser    func(*testing.T, *models.User)
	}{
		{
			name: "finds user by id successfully",
			id:   userID,
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch, userDoc))
			},
			wantErr: nil,
			checkUser: func(t *testing.T, u *models.User) {
				assert.Equal(t, userID, u.ID)
				assert.Equal(t, "alice", u.Username)
				assert.Equal(t, "alice@example.com", u.Email)
			},
		},
		{
			name: "returns error when user not found",
			id:   primitive.NewObjectID(),
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch))
			},
			wantErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			tt.addResponses(mt)

			svc := setupUserService(mt)
			user, err := svc.FindByID(context.Background(), tt.id)

			if tt.wantErr != nil {
				assert.ErrorIs(mt.T, err, tt.wantErr)
				assert.Nil(mt.T, user)
			} else {
				assert.NoError(mt.T, err)
				require.NotNil(mt.T, user)
				if tt.checkUser != nil {
					tt.checkUser(mt.T, user)
				}
			}
		})
	}
}

func TestUserService_FindByEmail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	userID := primitive.NewObjectID()
	now := time.Now()

	userDoc := bson.D{
		{Key: "_id", Value: userID},
		{Key: "username", Value: "alice"},
		{Key: "email", Value: "alice@example.com"},
		{Key: "password", Value: "$2a$10$hashedpassword"},
		{Key: "createdAt", Value: now},
		{Key: "updatedAt", Value: now},
	}

	tests := []struct {
		name         string
		email        string
		addResponses func(mt *mtest.T)
		wantErr      error
		checkUser    func(*testing.T, *models.User)
	}{
		{
			name:  "finds user by email successfully",
			email: "alice@example.com",
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch, userDoc))
			},
			wantErr: nil,
			checkUser: func(t *testing.T, u *models.User) {
				assert.Equal(t, "alice@example.com", u.Email)
				assert.Equal(t, "alice", u.Username)
			},
		},
		{
			name:  "returns error when email not found",
			email: "notfound@example.com",
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch))
			},
			wantErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			tt.addResponses(mt)

			svc := setupUserService(mt)
			user, err := svc.FindByEmail(context.Background(), tt.email)

			if tt.wantErr != nil {
				assert.ErrorIs(mt.T, err, tt.wantErr)
				assert.Nil(mt.T, user)
			} else {
				assert.NoError(mt.T, err)
				require.NotNil(mt.T, user)
				if tt.checkUser != nil {
					tt.checkUser(mt.T, user)
				}
			}
		})
	}
}

func TestUserService_Update(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	userID := primitive.NewObjectID()
	now := time.Now()

	newUsername := "alice_updated"
	newEmail := "new@example.com"

	updatedDoc := bson.D{
		{Key: "_id", Value: userID},
		{Key: "username", Value: "alice_updated"},
		{Key: "email", Value: "alice@example.com"},
		{Key: "password", Value: "$2a$10$hashedpassword"},
		{Key: "createdAt", Value: now},
		{Key: "updatedAt", Value: now},
	}

	tests := []struct {
		name         string
		id           primitive.ObjectID
		input        UpdateUserInput
		addResponses func(mt *mtest.T)
		wantErr      error
		checkUser    func(*testing.T, *models.User)
	}{
		{
			name: "updates username successfully",
			id:   userID,
			input: UpdateUserInput{
				Username: &newUsername,
			},
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(
					mtest.CreateSuccessResponse(),
					mtest.CreateCursorResponse(0, ns, mtest.FirstBatch, updatedDoc),
				)
			},
			wantErr: nil,
			checkUser: func(t *testing.T, u *models.User) {
				assert.Equal(t, "alice_updated", u.Username)
			},
		},
		{
			name: "returns error when updating to existing email",
			id:   userID,
			input: UpdateUserInput{
				Email: &newEmail,
			},
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(
					mtest.CreateCursorResponse(0, ns, mtest.FirstBatch, bson.D{{Key: "n", Value: int32(1)}}),
				)
			},
			wantErr: ErrEmailAlreadyExists,
		},
	}

	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			tt.addResponses(mt)

			svc := setupUserService(mt)
			user, err := svc.Update(context.Background(), tt.id, tt.input)

			if tt.wantErr != nil {
				assert.ErrorIs(mt.T, err, tt.wantErr)
				assert.Nil(mt.T, user)
			} else {
				assert.NoError(mt.T, err)
				require.NotNil(mt.T, user)
				if tt.checkUser != nil {
					tt.checkUser(mt.T, user)
				}
			}
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	tests := []struct {
		name         string
		id           primitive.ObjectID
		addResponses func(mt *mtest.T)
		wantErr      error
	}{
		{
			name: "deletes user successfully",
			id:   primitive.NewObjectID(),
			addResponses: func(mt *mtest.T) {
				mt.AddMockResponses(
					mtest.CreateSuccessResponse(bson.E{Key: "n", Value: 1}),
				)
			},
			wantErr: nil,
		},
		{
			name: "returns error when user not found",
			id:   primitive.NewObjectID(),
			addResponses: func(mt *mtest.T) {
				mt.AddMockResponses(
					mtest.CreateSuccessResponse(bson.E{Key: "n", Value: 0}),
				)
			},
			wantErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			tt.addResponses(mt)

			svc := setupUserService(mt)
			err := svc.Delete(context.Background(), tt.id)

			if tt.wantErr != nil {
				assert.ErrorIs(mt.T, err, tt.wantErr)
			} else {
				assert.NoError(mt.T, err)
			}
		})
	}
}

func TestUserService_Integration(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("full CRUD workflow", func(mt *mtest.T) {
		ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
		userID := primitive.NewObjectID()
		now := time.Now()

		userDoc := bson.D{
			{Key: "_id", Value: userID},
			{Key: "username", Value: "alice"},
			{Key: "email", Value: "alice@example.com"},
			{Key: "password", Value: "$2a$10$hashedpassword"},
			{Key: "createdAt", Value: now},
			{Key: "updatedAt", Value: now},
		}

		svc := setupUserService(mt)

		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, ns, mtest.FirstBatch),
			mtest.CreateSuccessResponse(),
		)
		user, err := svc.Create(context.Background(), CreateUserInput{
			Username: "alice",
			Email:    "alice@example.com",
			Password: "password123",
		})
		require.NoError(mt.T, err)
		require.NotNil(mt.T, user)
		assert.Equal(mt.T, "alice", user.Username)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch, userDoc))
		found, err := svc.FindByID(context.Background(), userID)
		require.NoError(mt.T, err)
		require.NotNil(mt.T, found)
		assert.Equal(mt.T, "alice@example.com", found.Email)

		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "n", Value: 1}))
		err = svc.Delete(context.Background(), userID)
		require.NoError(mt.T, err)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch))
		_, err = svc.FindByID(context.Background(), userID)
		assert.ErrorIs(mt.T, err, ErrUserNotFound)
	})
}
