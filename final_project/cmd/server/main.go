// Backend API
//
//	@title			Backend API
//	@version		1.0
//	@description	A RESTful API for managing users and their posts.
//
//	@host		localhost:8080
//	@BasePath	/api
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and your JWT token.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/Nick2603/golang/final_project/docs"

	"github.com/Nick2603/golang/final_project/cmd/server/handlers"
	"github.com/Nick2603/golang/final_project/cmd/server/middlewares"
	"github.com/Nick2603/golang/final_project/internal/clients"
	"github.com/Nick2603/golang/final_project/internal/config"
	"github.com/Nick2603/golang/final_project/internal/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/swagger"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	cfg := config.Load()

	db, err := clients.NewMongoDB(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		slog.Error("Failed to connect to MongoDB", "error", err)
		os.Exit(1)
	}
	defer db.Disconnect(context.Background())

	userSvc := services.NewUserService(db.Collection("users"))
	postSvc := services.NewPostService(db.Collection("posts"))

	authHandler := handlers.NewAuthHandler(userSvc, cfg.JWTSecret, cfg.JWTExpiry)
	userHandler := handlers.NewUserHandler(userSvc)
	postHandler := handlers.NewPostHandler(postSvc)

	app := fiber.New(fiber.Config{
		AppName: "Backend API v1.0",
	})

	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/signup", authHandler.SignUp)
	auth.Post("/signin", authHandler.SignIn)

	authMW := middlewares.Auth(cfg.JWTSecret)
	users := api.Group("/users", authMW)
	users.Get("/me", userHandler.GetMe)
	users.Patch("/me", userHandler.UpdateMe)
	users.Delete("/me", userHandler.DeleteMe)

	posts := api.Group("/posts")
	posts.Get("/", postHandler.GetAllPosts)
	posts.Get("/:id", postHandler.GetPostByID)
	posts.Post("/", authMW, postHandler.CreatePost)
	posts.Patch("/:id", authMW, postHandler.UpdatePost)
	posts.Delete("/:id", authMW, postHandler.DeletePost)

	go func() {
		slog.Info("Server starting", "port", cfg.Port)
		if err := app.Listen(":" + cfg.Port); err != nil {
			slog.Error("Server stopped with error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		slog.Error("Error during server shutdown", "error", err)
	}
	slog.Info("Server stopped")
}
