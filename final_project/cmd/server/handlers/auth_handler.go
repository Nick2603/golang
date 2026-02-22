package handlers

import (
	"log/slog"
	"time"

	httputils "github.com/Nick2603/golang/final_project/cmd/server/utils"
	"github.com/Nick2603/golang/final_project/internal/models"
	"github.com/Nick2603/golang/final_project/internal/services"
	"github.com/Nick2603/golang/final_project/internal/utils"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	userService *services.UserService
	jwtSecret   string
	jwtExpiry   time.Duration
}

func NewAuthHandler(userService *services.UserService, jwtSecret string, jwtExpiry time.Duration) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtSecret:   jwtSecret,
		jwtExpiry:   jwtExpiry,
	}
}

type SignUpRequest struct {
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email"    example:"john@example.com"`
	Password string `json:"password" example:"secret123"`
}

type SignInRequest struct {
	Email    string `json:"email"    example:"john@example.com"`
	Password string `json:"password" example:"secret123"`
}

type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// keep models in scope so swaggo resolves the type reference in annotations
var _ models.UserResponse

// SignUp godoc
//
//	@Summary		Register a new user
//	@Description	Create a new user account with username, email and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		SignUpRequest									true	"Sign-up credentials"
//	@Success		201		{object}	httputils.Response{data=models.UserResponse}	"User created"
//	@Failure		400		{object}	httputils.Response								"Invalid request"
//	@Failure		409		{object}	httputils.Response								"Email already taken"
//	@Failure		500		{object}	httputils.Response								"Internal error"
//	@Router			/auth/signup [post]
func (h *AuthHandler) SignUp(c fiber.Ctx) error {
	var req SignUpRequest
	if err := c.Bind().Body(&req); err != nil {
		slog.Warn("SignUp: failed to parse body", "error", err)
		return httputils.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return httputils.Error(c, fiber.StatusBadRequest, "username, email and password are required")
	}

	user, err := h.userService.Create(c.Context(), services.CreateUserInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch err {
		case services.ErrEmailAlreadyExists:
			return httputils.Error(c, fiber.StatusConflict, "email already exists")
		default:
			slog.Error("SignUp: failed to create user", "error", err)
			return httputils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}
	}

	slog.Info("SignUp: user registered", "userId", user.ID.Hex())
	return httputils.Success(c, fiber.StatusCreated, user.ToResponse())
}

// SignIn godoc
//
//	@Summary		Sign in
//	@Description	Authenticate with email and password; returns a JWT access token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		SignInRequest								true	"Sign-in credentials"
//	@Success		200		{object}	httputils.Response{data=TokenResponse}		"JWT token"
//	@Failure		400		{object}	httputils.Response							"Invalid request"
//	@Failure		401		{object}	httputils.Response							"Invalid credentials"
//	@Failure		500		{object}	httputils.Response							"Internal error"
//	@Router			/auth/signin [post]
func (h *AuthHandler) SignIn(c fiber.Ctx) error {
	var req SignInRequest
	if err := c.Bind().Body(&req); err != nil {
		slog.Warn("SignIn: failed to parse body", "error", err)
		return httputils.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	if req.Email == "" || req.Password == "" {
		return httputils.Error(c, fiber.StatusBadRequest, "email and password are required")
	}

	user, err := h.userService.FindByEmail(c.Context(), req.Email)
	if err != nil {
		slog.Warn("SignIn: user not found", "email", req.Email)
		return httputils.Error(c, fiber.StatusUnauthorized, "invalid credentials")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		slog.Warn("SignIn: wrong password", "email", req.Email)
		return httputils.Error(c, fiber.StatusUnauthorized, "invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID.Hex(), h.jwtSecret, h.jwtExpiry)
	if err != nil {
		slog.Error("SignIn: failed to generate token", "error", err)
		return httputils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	slog.Info("SignIn: user authenticated", "userId", user.ID.Hex())
	return httputils.Success(c, fiber.StatusOK, TokenResponse{Token: token})
}
