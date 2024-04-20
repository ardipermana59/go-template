package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardipermana59/go-template/config"
	"github.com/ardipermana59/go-template/internal/auth"
	"github.com/ardipermana59/go-template/internal/middleware"
	"github.com/ardipermana59/go-template/internal/post"
	"github.com/ardipermana59/go-template/internal/user"
	"github.com/ardipermana59/go-template/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	testDB         *gorm.DB
	testRouter     *gin.Engine
	testJWTService auth.JWTService
)

func setupTestDB(t *testing.T) {
	cfg, err := config.LoadConfig()
	assert.NoError(t, err)

	cfg.DBName = "testdb_test"
	
	db, err := database.NewDatabase(cfg.GetDSN())
	assert.NoError(t, err)

	db.Exec("DROP TABLE IF EXISTS posts")
	db.Exec("DROP TABLE IF EXISTS users")

	err = db.AutoMigrate(&user.User{}, &post.Post{})
	assert.NoError(t, err)

	testDB = db
	testJWTService = auth.NewJWTService(cfg.JWTSecret, cfg.JWTExpireHours)
}

func setupTestRouter() {
	userRepo := user.NewRepository(testDB)
	userService := user.NewService(userRepo, testJWTService)
	userHandler := user.NewHandler(userService)

	postRepo := post.NewRepository(testDB)
	postService := post.NewService(postRepo)
	postHandler := post.NewHandler(postService)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", userHandler.Register)
			authGroup.POST("/login", userHandler.Login)
		}

		protectedGroup := api.Group("")
		protectedGroup.Use(middleware.AuthMiddleware(testJWTService))
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
		adminGroup.Use(middleware.AuthMiddleware(testJWTService))
		adminGroup.Use(middleware.RoleMiddleware("admin"))
		{
			adminGroup.GET("/users", userHandler.GetAllUsers)
			adminGroup.GET("/users/:id", userHandler.GetUserByID)
			adminGroup.PUT("/users/:id", userHandler.UpdateUser)
			adminGroup.DELETE("/users/:id", userHandler.DeleteUser)
		}
	}

	testRouter = r
}

func TestRegisterUser(t *testing.T) {
	setupTestDB(t)
	setupTestRouter()

	t.Run("Success - Register new user", func(t *testing.T) {
		payload := map[string]string{
			"name":             "John Doe",
			"email":            "john@example.com",
			"password":         "password123",
			"password_confirm": "password123",
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, true, response["success"])
		assert.Equal(t, "User registered successfully", response["message"])
	})

	t.Run("Fail - Email already exists", func(t *testing.T) {
		payload := map[string]string{
			"name":             "Jane Doe",
			"email":            "john@example.com",
			"password":         "password123",
			"password_confirm": "password123",
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, false, response["success"])
		assert.Equal(t, "Failed to register user", response["message"])
		
		errors := response["error"].([]interface{})
		firstError := errors[0].(map[string]interface{})
		assert.Equal(t, "email", firstError["field"])
		assert.Equal(t, "The email has already been taken", firstError["message"])
	})

	t.Run("Fail - Validation errors", func(t *testing.T) {
		payload := map[string]string{
			"name":             "Jo",
			"email":            "invalid-email",
			"password":         "12345",
			"password_confirm": "123456",
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, false, response["success"])
		assert.Equal(t, "Validation failed", response["message"])
		assert.NotNil(t, response["error"])
	})
}

func TestLoginUser(t *testing.T) {
	setupTestDB(t)
	setupTestRouter()

	registerPayload := map[string]string{
		"name":             "Test User",
		"email":            "test@example.com",
		"password":         "password123",
		"password_confirm": "password123",
	}
	body, _ := json.Marshal(registerPayload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	t.Run("Success - Login with valid credentials", func(t *testing.T) {
		payload := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, true, response["success"])
		assert.Equal(t, "Login successful", response["message"])
		
		data := response["data"].(map[string]interface{})
		assert.NotEmpty(t, data["token"])
	})

	t.Run("Fail - Invalid credentials", func(t *testing.T) {
		payload := map[string]string{
			"email":    "test@example.com",
			"password": "wrongpassword",
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, false, response["success"])
		
		errors := response["error"].([]interface{})
		firstError := errors[0].(map[string]interface{})
		assert.Equal(t, "credentials", firstError["field"])
		assert.Equal(t, "The provided credentials are invalid", firstError["message"])
	})
}

func TestGetProfile(t *testing.T) {
	setupTestDB(t)
	setupTestRouter()

	registerPayload := map[string]string{
		"name":             "Profile User",
		"email":            "profile@example.com",
		"password":         "password123",
		"password_confirm": "password123",
	}
	body, _ := json.Marshal(registerPayload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	loginPayload := map[string]string{
		"email":    "profile@example.com",
		"password": "password123",
	}
	body, _ = json.Marshal(loginPayload)
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	data := loginResponse["data"].(map[string]interface{})
	token := data["token"].(string)

	t.Run("Success - Get profile with valid token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, true, response["success"])
	})

	t.Run("Fail - No authorization header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/profile", nil)

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, false, response["success"])
	})
}
