package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Nick2603/golang/final_project/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrForbidden    = errors.New("forbidden: resource belongs to another user")
)

// PostService handles business logic for posts.
type PostService struct {
	col *mongo.Collection
}

func NewPostService(col *mongo.Collection) *PostService {
	return &PostService{col: col}
}

// CreatePostInput contains data required to create a post.
type CreatePostInput struct {
	UserID  primitive.ObjectID
	Title   string
	Content string
}

// UpdatePostInput contains optional fields to update on a post.
type UpdatePostInput struct {
	Title   *string
	Content *string
}

// Create inserts a new post.
func (s *PostService) Create(ctx context.Context, input CreatePostInput) (*models.Post, error) {
	now := time.Now()
	post := &models.Post{
		ID:        primitive.NewObjectID(),
		UserID:    input.UserID,
		Title:     input.Title,
		Content:   input.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if _, err := s.col.InsertOne(ctx, post); err != nil {
		return nil, err
	}

	slog.Info("Post created", "id", post.ID.Hex(), "userId", input.UserID.Hex())
	return post, nil
}

// FindAll returns all posts sorted by creation date descending.
func (s *PostService) FindAll(ctx context.Context) ([]*models.Post, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := s.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*models.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	slog.Info("Posts fetched", "count", len(posts))
	return posts, nil
}

// FindByID retrieves a single post by its ObjectID.
func (s *PostService) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Post, error) {
	var post models.Post
	err := s.col.FindOne(ctx, bson.M{"_id": id}).Decode(&post)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}
	slog.Info("Post fetched by ID", "id", id.Hex())
	return &post, nil
}

// Update applies a partial update to a post, enforcing ownership.
func (s *PostService) Update(ctx context.Context, id, userID primitive.ObjectID, input UpdatePostInput) (*models.Post, error) {
	post, err := s.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if post.UserID != userID {
		return nil, ErrForbidden
	}

	fields := bson.M{"updatedAt": time.Now()}
	if input.Title != nil {
		fields["title"] = *input.Title
	}
	if input.Content != nil {
		fields["content"] = *input.Content
	}

	if _, err := s.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": fields}); err != nil {
		return nil, err
	}

	slog.Info("Post updated", "id", id.Hex())
	return s.FindByID(ctx, id)
}

// Delete removes a post, enforcing ownership.
func (s *PostService) Delete(ctx context.Context, id, userID primitive.ObjectID) error {
	post, err := s.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if post.UserID != userID {
		return ErrForbidden
	}

	if _, err := s.col.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return err
	}

	slog.Info("Post deleted", "id", id.Hex())
	return nil
}
