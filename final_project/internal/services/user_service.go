package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Nick2603/golang/final_project/internal/models"
	"github.com/Nick2603/golang/final_project/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type UserService struct {
	col *mongo.Collection
}

func NewUserService(col *mongo.Collection) *UserService {
	return &UserService{col: col}
}

type CreateUserInput struct {
	Username string
	Email    string
	Password string
}

type UpdateUserInput struct {
	Username *string
	Email    *string
	Password *string
}

func (s *UserService) Create(ctx context.Context, input CreateUserInput) (*models.User, error) {
	count, err := s.col.CountDocuments(ctx, bson.M{"email": input.Email})
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrEmailAlreadyExists
	}

	hashedPw, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &models.User{
		ID:        primitive.NewObjectID(),
		Username:  input.Username,
		Email:     input.Email,
		Password:  hashedPw,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if _, err := s.col.InsertOne(ctx, user); err != nil {
		return nil, err
	}

	slog.Info("User created", "id", user.ID.Hex(), "email", user.Email)
	return user, nil
}

func (s *UserService) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := s.col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	slog.Info("User fetched by ID", "id", id.Hex())
	return &user, nil
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := s.col.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Update(ctx context.Context, id primitive.ObjectID, input UpdateUserInput) (*models.User, error) {
	fields := bson.M{"updatedAt": time.Now()}

	if input.Username != nil {
		fields["username"] = *input.Username
	}
	if input.Email != nil {
		count, err := s.col.CountDocuments(ctx, bson.M{
			"email": *input.Email,
			"_id":   bson.M{"$ne": id},
		})
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, ErrEmailAlreadyExists
		}
		fields["email"] = *input.Email
	}
	if input.Password != nil {
		hashed, err := utils.HashPassword(*input.Password)
		if err != nil {
			return nil, err
		}
		fields["password"] = hashed
	}

	if _, err := s.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": fields}); err != nil {
		return nil, err
	}

	slog.Info("User updated", "id", id.Hex())
	return s.FindByID(ctx, id)
}

func (s *UserService) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := s.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrUserNotFound
	}
	slog.Info("User deleted", "id", id.Hex())
	return nil
}
