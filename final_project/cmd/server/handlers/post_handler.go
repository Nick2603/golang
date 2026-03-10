package handlers

import (
	"log/slog"

	httputils "github.com/Nick2603/golang/final_project/cmd/server/utils"
	"github.com/Nick2603/golang/final_project/internal/models"
	"github.com/Nick2603/golang/final_project/internal/services"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

type CreatePostRequest struct {
	Title   string `json:"title"   example:"My first post"`
	Content string `json:"content" example:"Hello world!"`
}

type UpdatePostRequest struct {
	Title   *string `json:"title"   example:"Updated title"`
	Content *string `json:"content" example:"Updated content"`
}

// keep models in scope for swaggo
var _ models.Post

// CreatePost godoc
//
//	@Summary		Create a post
//	@Description	Creates a new post for the authenticated user
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		CreatePostRequest						true	"Post data"
//	@Success		201		{object}	httputils.Response{data=models.Post}	"Post created"
//	@Failure		400		{object}	httputils.Response						"Invalid request"
//	@Failure		401		{object}	httputils.Response						"Unauthorized"
//	@Failure		500		{object}	httputils.Response						"Internal error"
//	@Router			/posts [post]
func (h *PostHandler) CreatePost(c fiber.Ctx) error {
	userID, err := userIDFromCtx(c)
	if err != nil {
		return httputils.Error(c, fiber.StatusBadRequest, "invalid user id in token")
	}

	var req CreatePostRequest
	if err := c.Bind().Body(&req); err != nil {
		slog.Warn("CreatePost: failed to parse body", "error", err)
		return httputils.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	if req.Title == "" || req.Content == "" {
		return httputils.Error(c, fiber.StatusBadRequest, "title and content are required")
	}

	post, err := h.postService.Create(c.Context(), services.CreatePostInput{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		slog.Error("CreatePost: failed to create post", "error", err, "userID", userID.Hex())
		return httputils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	return httputils.Success(c, fiber.StatusCreated, post)
}

// GetAllPosts godoc
//
//	@Summary		Get all posts
//	@Description	Returns all posts sorted by newest first (public endpoint)
//	@Tags			posts
//	@Produce		json
//	@Success		200	{object}	httputils.Response{data=[]models.Post}	"List of posts"
//	@Failure		500	{object}	httputils.Response						"Internal error"
//	@Router			/posts [get]
func (h *PostHandler) GetAllPosts(c fiber.Ctx) error {
	posts, err := h.postService.FindAll(c.Context())
	if err != nil {
		slog.Error("GetAllPosts: failed to fetch posts", "error", err)
		return httputils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	// Return empty array instead of null when no posts exist.
	if posts == nil {
		posts = []*models.Post{}
	}

	return httputils.Success(c, fiber.StatusOK, posts)
}

// GetPostByID godoc
//
//	@Summary		Get post by ID
//	@Description	Returns a single post by its ID (public endpoint)
//	@Tags			posts
//	@Produce		json
//	@Param			id	path		string									true	"Post ID (ObjectID hex)"
//	@Success		200	{object}	httputils.Response{data=models.Post}	"Post"
//	@Failure		400	{object}	httputils.Response						"Invalid ID"
//	@Failure		404	{object}	httputils.Response						"Not found"
//	@Failure		500	{object}	httputils.Response						"Internal error"
//	@Router			/posts/{id} [get]
func (h *PostHandler) GetPostByID(c fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return httputils.Error(c, fiber.StatusBadRequest, "invalid post id")
	}

	post, err := h.postService.FindByID(c.Context(), id)
	if err != nil {
		switch err {
		case services.ErrPostNotFound:
			return httputils.Error(c, fiber.StatusNotFound, "post not found")
		default:
			slog.Error("GetPostByID: failed to fetch post", "error", err, "postID", id.Hex())
			return httputils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}
	}

	return httputils.Success(c, fiber.StatusOK, post)
}

// UpdatePost godoc
//
//	@Summary		Update a post
//	@Description	Partially updates a post; only the owner can update
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string									true	"Post ID (ObjectID hex)"
//	@Param			request	body		UpdatePostRequest						true	"Fields to update (all optional)"
//	@Success		200		{object}	httputils.Response{data=models.Post}	"Updated post"
//	@Failure		400		{object}	httputils.Response						"Invalid request"
//	@Failure		401		{object}	httputils.Response						"Unauthorized"
//	@Failure		403		{object}	httputils.Response						"Forbidden"
//	@Failure		404		{object}	httputils.Response						"Not found"
//	@Failure		500		{object}	httputils.Response						"Internal error"
//	@Router			/posts/{id} [patch]
func (h *PostHandler) UpdatePost(c fiber.Ctx) error {
	userID, err := userIDFromCtx(c)
	if err != nil {
		return httputils.Error(c, fiber.StatusBadRequest, "invalid user id in token")
	}

	postID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return httputils.Error(c, fiber.StatusBadRequest, "invalid post id")
	}

	var req UpdatePostRequest
	if err := c.Bind().Body(&req); err != nil {
		slog.Warn("UpdatePost: failed to parse body", "error", err)
		return httputils.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	post, err := h.postService.Update(c.Context(), postID, userID, services.UpdatePostInput{
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		switch err {
		case services.ErrPostNotFound:
			return httputils.Error(c, fiber.StatusNotFound, "post not found")
		case services.ErrForbidden:
			return httputils.Error(c, fiber.StatusForbidden, "you can only update your own posts")
		default:
			slog.Error("UpdatePost: failed to update post", "error", err, "postID", postID.Hex())
			return httputils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}
	}

	return httputils.Success(c, fiber.StatusOK, post)
}

// DeletePost godoc
//
//	@Summary		Delete a post
//	@Description	Deletes a post; only the owner can delete
//	@Tags			posts
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string				true	"Post ID (ObjectID hex)"
//	@Success		204	"No content"
//	@Failure		400	{object}	httputils.Response	"Invalid ID"
//	@Failure		401	{object}	httputils.Response	"Unauthorized"
//	@Failure		403	{object}	httputils.Response	"Forbidden"
//	@Failure		404	{object}	httputils.Response	"Not found"
//	@Failure		500	{object}	httputils.Response	"Internal error"
//	@Router			/posts/{id} [delete]
func (h *PostHandler) DeletePost(c fiber.Ctx) error {
	userID, err := userIDFromCtx(c)
	if err != nil {
		return httputils.Error(c, fiber.StatusBadRequest, "invalid user id in token")
	}

	postID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return httputils.Error(c, fiber.StatusBadRequest, "invalid post id")
	}

	if err := h.postService.Delete(c.Context(), postID, userID); err != nil {
		switch err {
		case services.ErrPostNotFound:
			return httputils.Error(c, fiber.StatusNotFound, "post not found")
		case services.ErrForbidden:
			return httputils.Error(c, fiber.StatusForbidden, "you can only delete your own posts")
		default:
			slog.Error("DeletePost: failed to delete post", "error", err, "postID", postID.Hex())
			return httputils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}
