package post

import (
	"net/http"
	"strconv"

	"github.com/ardipermana59/go-template/internal/common/apperror"
	"github.com/ardipermana59/go-template/internal/common/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreatePost(c *gin.Context) {
	userID := c.GetUint("user_id")

	var dto CreatePostDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, err)
		return
	}

	post, appErr := h.service.CreatePost(userID, dto)
	if appErr != nil {
		response.InternalError(c, nil)
		return
	}

	response.Success(c, http.StatusCreated, "Post created successfully", post)
}

func (h *Handler) GetAllPosts(c *gin.Context) {
	posts, appErr := h.service.GetAllPosts()
	if appErr != nil {
		response.InternalError(c, nil)
		return
	}

	response.Success(c, http.StatusOK, "Posts retrieved successfully", posts)
}

func (h *Handler) GetPostByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", apperror.InvalidID())
		return
	}

	post, appErr := h.service.GetPostByID(uint(id))
	if appErr != nil {
		response.Error(c, http.StatusNotFound, "Post not found", appErr)
		return
	}

	response.Success(c, http.StatusOK, "Post retrieved successfully", post)
}

func (h *Handler) GetPostsByUserID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID", apperror.InvalidID())
		return
	}

	posts, appErr := h.service.GetPostsByUserID(uint(userID))
	if appErr != nil {
		response.InternalError(c, nil)
		return
	}

	response.Success(c, http.StatusOK, "Posts retrieved successfully", posts)
}

func (h *Handler) GetMyPosts(c *gin.Context) {
	userID := c.GetUint("user_id")

	posts, appErr := h.service.GetMyPosts(userID)
	if appErr != nil {
		response.InternalError(c, nil)
		return
	}

	response.Success(c, http.StatusOK, "Posts retrieved successfully", posts)
}

func (h *Handler) UpdatePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", apperror.InvalidID())
		return
	}

	var dto UpdatePostDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, err)
		return
	}

	post, appErr := h.service.UpdatePost(uint(id), userID, dto)
	if appErr != nil {
		response.Error(c, http.StatusForbidden, "Failed to update post", appErr)
		return
	}

	response.Success(c, http.StatusOK, "Post updated successfully", post)
}

func (h *Handler) DeletePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", apperror.InvalidID())
		return
	}

	appErr := h.service.DeletePost(uint(id), userID)
	if appErr != nil {
		response.Error(c, http.StatusForbidden, "Failed to delete post", appErr)
		return
	}

	response.Success(c, http.StatusOK, "Post deleted successfully", nil)
}
