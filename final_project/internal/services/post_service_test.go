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

func setupPostService(mt *mtest.T) *PostService {
	mt.Helper()
	return NewPostService(mt.Coll)
}

func TestPostService_Create(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	userID := primitive.NewObjectID()

	tests := []struct {
		name         string
		input        CreatePostInput
		addResponses func(mt *mtest.T)
		wantErr      error
		checkPost    func(*testing.T, *models.Post)
	}{
		{
			name: "creates post successfully",
			input: CreatePostInput{
				UserID:  userID,
				Title:   "My first post",
				Content: "Hello world!",
			},
			addResponses: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: nil,
			checkPost: func(t *testing.T, p *models.Post) {
				assert.Equal(t, userID, p.UserID)
				assert.Equal(t, "My first post", p.Title)
				assert.Equal(t, "Hello world!", p.Content)
				assert.False(t, p.ID.IsZero())
			},
		},
		{
			name: "creates post with different content",
			input: CreatePostInput{
				UserID:  userID,
				Title:   "Second post",
				Content: "More content here",
			},
			addResponses: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: nil,
			checkPost: func(t *testing.T, p *models.Post) {
				assert.Equal(t, "Second post", p.Title)
				assert.Equal(t, "More content here", p.Content)
			},
		},
	}

	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			tt.addResponses(mt)

			svc := setupPostService(mt)
			post, err := svc.Create(context.Background(), tt.input)

			if tt.wantErr != nil {
				assert.ErrorIs(mt.T, err, tt.wantErr)
				assert.Nil(mt.T, post)
			} else {
				assert.NoError(mt.T, err)
				require.NotNil(mt.T, post)
				if tt.checkPost != nil {
					tt.checkPost(mt.T, post)
				}
			}
		})
	}
}

func TestPostService_FindAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	userID := primitive.NewObjectID()
	now := time.Now()

	post1Doc := bson.D{
		{Key: "_id", Value: primitive.NewObjectID()},
		{Key: "userId", Value: userID},
		{Key: "title", Value: "First post"},
		{Key: "content", Value: "Content 1"},
		{Key: "createdAt", Value: now},
		{Key: "updatedAt", Value: now},
	}
	post2Doc := bson.D{
		{Key: "_id", Value: primitive.NewObjectID()},
		{Key: "userId", Value: userID},
		{Key: "title", Value: "Second post"},
		{Key: "content", Value: "Content 2"},
		{Key: "createdAt", Value: now},
		{Key: "updatedAt", Value: now},
	}

	tests := []struct {
		name         string
		addResponses func(mt *mtest.T)
		wantErr      error
		wantCount    int
		checkPosts   func(*testing.T, []*models.Post)
	}{
		{
			name: "returns all posts",
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch, post1Doc, post2Doc))
			},
			wantErr:   nil,
			wantCount: 2,
			checkPosts: func(t *testing.T, posts []*models.Post) {
				assert.Equal(t, "First post", posts[0].Title)
				assert.Equal(t, "Second post", posts[1].Title)
			},
		},
		{
			name: "returns empty list when no posts exist",
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch))
			},
			wantErr:   nil,
			wantCount: 0,
			checkPosts: func(t *testing.T, posts []*models.Post) {
				assert.Empty(t, posts)
			},
		},
	}

	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			tt.addResponses(mt)

			svc := setupPostService(mt)
			posts, err := svc.FindAll(context.Background())

			if tt.wantErr != nil {
				assert.ErrorIs(mt.T, err, tt.wantErr)
				assert.Nil(mt.T, posts)
			} else {
				assert.NoError(mt.T, err)
				assert.Len(mt.T, posts, tt.wantCount)
				if tt.checkPosts != nil {
					tt.checkPosts(mt.T, posts)
				}
			}
		})
	}
}

func TestPostService_FindByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	postID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	now := time.Now()

	postDoc := bson.D{
		{Key: "_id", Value: postID},
		{Key: "userId", Value: userID},
		{Key: "title", Value: "My post"},
		{Key: "content", Value: "Post content"},
		{Key: "createdAt", Value: now},
		{Key: "updatedAt", Value: now},
	}

	tests := []struct {
		name         string
		id           primitive.ObjectID
		addResponses func(mt *mtest.T)
		wantErr      error
		checkPost    func(*testing.T, *models.Post)
	}{
		{
			name: "finds post by id successfully",
			id:   postID,
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch, postDoc))
			},
			wantErr: nil,
			checkPost: func(t *testing.T, p *models.Post) {
				assert.Equal(t, postID, p.ID)
				assert.Equal(t, userID, p.UserID)
				assert.Equal(t, "My post", p.Title)
				assert.Equal(t, "Post content", p.Content)
			},
		},
		{
			name: "returns error when post not found",
			id:   primitive.NewObjectID(),
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch))
			},
			wantErr: ErrPostNotFound,
		},
	}

	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			tt.addResponses(mt)

			svc := setupPostService(mt)
			post, err := svc.FindByID(context.Background(), tt.id)

			if tt.wantErr != nil {
				assert.ErrorIs(mt.T, err, tt.wantErr)
				assert.Nil(mt.T, post)
			} else {
				assert.NoError(mt.T, err)
				require.NotNil(mt.T, post)
				if tt.checkPost != nil {
					tt.checkPost(mt.T, post)
				}
			}
		})
	}
}

func TestPostService_Update(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	postID := primitive.NewObjectID()
	ownerID := primitive.NewObjectID()
	otherUserID := primitive.NewObjectID()
	now := time.Now()

	newTitle := "Updated title"

	postDoc := bson.D{
		{Key: "_id", Value: postID},
		{Key: "userId", Value: ownerID},
		{Key: "title", Value: "Original title"},
		{Key: "content", Value: "Original content"},
		{Key: "createdAt", Value: now},
		{Key: "updatedAt", Value: now},
	}
	updatedDoc := bson.D{
		{Key: "_id", Value: postID},
		{Key: "userId", Value: ownerID},
		{Key: "title", Value: "Updated title"},
		{Key: "content", Value: "Original content"},
		{Key: "createdAt", Value: now},
		{Key: "updatedAt", Value: now},
	}

	tests := []struct {
		name         string
		postID       primitive.ObjectID
		userID       primitive.ObjectID
		input        UpdatePostInput
		addResponses func(mt *mtest.T)
		wantErr      error
		checkPost    func(*testing.T, *models.Post)
	}{
		{
			name:   "updates post successfully",
			postID: postID,
			userID: ownerID,
			input:  UpdatePostInput{Title: &newTitle},
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(
					mtest.CreateCursorResponse(0, ns, mtest.FirstBatch, postDoc),
					mtest.CreateSuccessResponse(),
					mtest.CreateCursorResponse(0, ns, mtest.FirstBatch, updatedDoc),
				)
			},
			wantErr: nil,
			checkPost: func(t *testing.T, p *models.Post) {
				assert.Equal(t, "Updated title", p.Title)
			},
		},
		{
			name:   "returns forbidden when user is not the owner",
			postID: postID,
			userID: otherUserID,
			input:  UpdatePostInput{Title: &newTitle},
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(
					mtest.CreateCursorResponse(0, ns, mtest.FirstBatch, postDoc),
				)
			},
			wantErr: ErrForbidden,
		},
		{
			name:   "returns error when post not found",
			postID: primitive.NewObjectID(),
			userID: ownerID,
			input:  UpdatePostInput{Title: &newTitle},
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(
					mtest.CreateCursorResponse(0, ns, mtest.FirstBatch),
				)
			},
			wantErr: ErrPostNotFound,
		},
	}

	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			tt.addResponses(mt)

			svc := setupPostService(mt)
			post, err := svc.Update(context.Background(), tt.postID, tt.userID, tt.input)

			if tt.wantErr != nil {
				assert.ErrorIs(mt.T, err, tt.wantErr)
				assert.Nil(mt.T, post)
			} else {
				assert.NoError(mt.T, err)
				require.NotNil(mt.T, post)
				if tt.checkPost != nil {
					tt.checkPost(mt.T, post)
				}
			}
		})
	}
}

func TestPostService_Delete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	postID := primitive.NewObjectID()
	ownerID := primitive.NewObjectID()
	otherUserID := primitive.NewObjectID()
	now := time.Now()

	postDoc := bson.D{
		{Key: "_id", Value: postID},
		{Key: "userId", Value: ownerID},
		{Key: "title", Value: "My post"},
		{Key: "content", Value: "Content"},
		{Key: "createdAt", Value: now},
		{Key: "updatedAt", Value: now},
	}

	tests := []struct {
		name         string
		postID       primitive.ObjectID
		userID       primitive.ObjectID
		addResponses func(mt *mtest.T)
		wantErr      error
	}{
		{
			name:   "deletes post successfully",
			postID: postID,
			userID: ownerID,
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(
					mtest.CreateCursorResponse(0, ns, mtest.FirstBatch, postDoc),
					mtest.CreateSuccessResponse(),
				)
			},
			wantErr: nil,
		},
		{
			name:   "returns forbidden when user is not the owner",
			postID: postID,
			userID: otherUserID,
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(
					mtest.CreateCursorResponse(0, ns, mtest.FirstBatch, postDoc),
				)
			},
			wantErr: ErrForbidden,
		},
		{
			name:   "returns error when post not found",
			postID: primitive.NewObjectID(),
			userID: ownerID,
			addResponses: func(mt *mtest.T) {
				ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
				mt.AddMockResponses(
					mtest.CreateCursorResponse(0, ns, mtest.FirstBatch),
				)
			},
			wantErr: ErrPostNotFound,
		},
	}

	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			tt.addResponses(mt)

			svc := setupPostService(mt)
			err := svc.Delete(context.Background(), tt.postID, tt.userID)

			if tt.wantErr != nil {
				assert.ErrorIs(mt.T, err, tt.wantErr)
			} else {
				assert.NoError(mt.T, err)
			}
		})
	}
}

func TestPostService_Integration(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("full CRUD workflow", func(mt *mtest.T) {
		ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
		postID := primitive.NewObjectID()
		userID := primitive.NewObjectID()
		now := time.Now()

		postDoc := bson.D{
			{Key: "_id", Value: postID},
			{Key: "userId", Value: userID},
			{Key: "title", Value: "My post"},
			{Key: "content", Value: "Some content"},
			{Key: "createdAt", Value: now},
			{Key: "updatedAt", Value: now},
		}

		svc := setupPostService(mt)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		post, err := svc.Create(context.Background(), CreatePostInput{
			UserID:  userID,
			Title:   "My post",
			Content: "Some content",
		})
		require.NoError(mt.T, err)
		require.NotNil(mt.T, post)
		assert.Equal(mt.T, "My post", post.Title)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch, postDoc))
		found, err := svc.FindByID(context.Background(), postID)
		require.NoError(mt.T, err)
		require.NotNil(mt.T, found)
		assert.Equal(mt.T, "Some content", found.Content)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch, postDoc))
		posts, err := svc.FindAll(context.Background())
		require.NoError(mt.T, err)
		assert.Len(mt.T, posts, 1)

		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, ns, mtest.FirstBatch, postDoc),
			mtest.CreateSuccessResponse(),
		)
		err = svc.Delete(context.Background(), postID, userID)
		require.NoError(mt.T, err)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch))
		_, err = svc.FindByID(context.Background(), postID)
		assert.ErrorIs(mt.T, err, ErrPostNotFound)
	})
}
