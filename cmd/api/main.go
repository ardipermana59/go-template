package main

import (
	"log"

	"github.com/ardipermana59/go-template/config"
	"github.com/ardipermana59/go-template/internal/auth"
	"github.com/ardipermana59/go-template/internal/middleware"
	"github.com/ardipermana59/go-template/internal/post"
	"github.com/ardipermana59/go-template/internal/user"
	"github.com/ardipermana59/go-template/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := database.NewDatabase(cfg.GetDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&user.User{}, &post.Post{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.JWTExpireHours)

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, jwtService)
	userHandler := user.NewHandler(userService)

	postRepo := post.NewRepository(db)
	postService := post.NewService(postRepo)
	postHandler := post.NewHandler(postService)

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", userHandler.Register)
			authGroup.POST("/login", userHandler.Login)
		}

		protectedGroup := api.Group("")
		protectedGroup.Use(middleware.AuthMiddleware(jwtService))
		{
			protectedGroup.GET("/profile", userHandler.GetProfile)
			protectedGroup.PUT("/profile", userHandler.UpdateProfile)
			protectedGroup.PUT("/change-password", userHandler.ChangePassword)

			protectedGroup.GET("/posts/my", postHandler.GetMyPosts)
			protectedGroup.POST("/posts", postHandler.CreatePost)
			protectedGroup.PUT("/posts/:id", postHandler.UpdatePost)
			protectedGroup.DELETE("/posts/:id", postHandler.DeletePost)
		}

		publicGroup := api.Group("")
		{
			publicGroup.GET("/posts", postHandler.GetAllPosts)
			publicGroup.GET("/posts/:id", postHandler.GetPostByID)
			publicGroup.GET("/users/:user_id/posts", postHandler.GetPostsByUserID)
		}

		adminGroup := api.Group("/admin")
		adminGroup.Use(middleware.AuthMiddleware(jwtService))
		adminGroup.Use(middleware.RoleMiddleware("admin"))
		{
			adminGroup.GET("/users", userHandler.GetAllUsers)
			adminGroup.GET("/users/:id", userHandler.GetUserByID)
			adminGroup.PUT("/users/:id", userHandler.UpdateUser)
			adminGroup.DELETE("/users/:id", userHandler.DeleteUser)
		}
	}

	log.Printf("🚀 Server running on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
