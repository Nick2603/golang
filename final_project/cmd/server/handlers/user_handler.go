package handlers

import (
	"log/slog"

	httputils "github.com/Nick2603/golang/final_project/cmd/server/utils"
	"github.com/Nick2603/golang/final_project/internal/models"
	"github.com/Nick2603/golang/final_project/internal/services"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type UpdateUserRequest struct {
	Username *string `json:"username" example:"newname"`
	Email    *string `json:"email"    example:"new@example.com"`
	Password *string `json:"password" example:"newpassword"`
}

// keep models in scope for swaggo
var _ models.UserResponse

// GetMe godoc
//
//	@Summary		Get current user
//	@Description	Returns the profile of the authenticated user
//	@Tags			users
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	httputils.Response{data=models.UserResponse}	"User profile"
//	@Failure		401	{object}	httputils.Response								"Unauthorized"
//	@Failure		404	{object}	httputils.Response								"User not found"
//	@Failure		500	{object}	httputils.Response								"Internal error"
//	@Router			/users/me [get]
func (h *UserHandler) GetMe(c fiber.Ctx) error {
	userID, err := userIDFromCtx(c)
	if err != nil {
		return httputils.Error(c, fiber.StatusBadRequest, "invalid user id in token")
	}

	user, err := h.userService.FindByID(c.Context(), userID)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			return httputils.Error(c, fiber.StatusNotFound, "user not found")
		default:
			slog.Error("GetMe: failed to find user", "error", err, "userID", userID.Hex())
			return httputils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}
	}

	return httputils.Success(c, fiber.StatusOK, user.ToResponse())
}

// UpdateMe godoc
//
//	@Summary		Update current user
//	@Description	Partially updates the authenticated user's profile
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		UpdateUserRequest								true	"Fields to update (all optional)"
//	@Success		200		{object}	httputils.Response{data=models.UserResponse}	"Updated user"
//	@Failure		400		{object}	httputils.Response								"Invalid request"
//	@Failure		401		{object}	httputils.Response								"Unauthorized"
//	@Failure		409		{object}	httputils.Response								"Email already taken"
//	@Failure		500		{object}	httputils.Response								"Internal error"
//	@Router			/users/me [patch]
func (h *UserHandler) UpdateMe(c fiber.Ctx) error {
	userID, err := userIDFromCtx(c)
	if err != nil {
		return httputils.Error(c, fiber.StatusBadRequest, "invalid user id in token")
	}

	var req UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		slog.Warn("UpdateMe: failed to parse body", "error", err)
		return httputils.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	user, err := h.userService.Update(c.Context(), userID, services.UpdateUserInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch err {
		case services.ErrEmailAlreadyExists:
			return httputils.Error(c, fiber.StatusConflict, "email already exists")
		default:
			slog.Error("UpdateMe: failed to update user", "error", err, "userID", userID.Hex())
			return httputils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}
	}

	return httputils.Success(c, fiber.StatusOK, user.ToResponse())
}

// DeleteMe godoc
//
//	@Summary		Delete current user
//	@Description	Permanently deletes the authenticated user's account
//	@Tags			users
//	@Produce		json
//	@Security		BearerAuth
//	@Success		204	"No content"
//	@Failure		401	{object}	httputils.Response	"Unauthorized"
//	@Failure		500	{object}	httputils.Response	"Internal error"
//	@Router			/users/me [delete]
func (h *UserHandler) DeleteMe(c fiber.Ctx) error {
	userID, err := userIDFromCtx(c)
	if err != nil {
		return httputils.Error(c, fiber.StatusBadRequest, "invalid user id in token")
	}

	if err := h.userService.Delete(c.Context(), userID); err != nil {
		slog.Error("DeleteMe: failed to delete user", "error", err, "userID", userID.Hex())
		return httputils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	slog.Info("DeleteMe: user deleted", "userID", userID.Hex())
	return c.SendStatus(fiber.StatusNoContent)
}

func userIDFromCtx(c fiber.Ctx) (primitive.ObjectID, error) {
	raw, _ := c.Locals("userID").(string)
	return primitive.ObjectIDFromHex(raw)
}
