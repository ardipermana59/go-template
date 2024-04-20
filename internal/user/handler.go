package user

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

func (h *Handler) Register(c *gin.Context) {
	var dto RegisterDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, err)
		return
	}

	user, appErr := h.service.Register(dto)
	if appErr != nil {
		response.Error(c, http.StatusBadRequest, "Failed to register user", appErr)
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", user)
}

func (h *Handler) Login(c *gin.Context) {
	var dto LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, err)
		return
	}

	result, appErr := h.service.Login(dto)
	if appErr != nil {
		response.Error(c, http.StatusUnauthorized, "Login failed", appErr)
		return
	}

	response.Success(c, http.StatusOK, "Login successful", result)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, appErr := h.service.GetProfile(userID)
	if appErr != nil {
		response.Error(c, http.StatusNotFound, "User not found", appErr)
		return
	}

	response.Success(c, http.StatusOK, "Profile retrieved successfully", user)
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	users, appErr := h.service.GetAllUsers()
	if appErr != nil {
		response.InternalError(c, nil)
		return
	}

	response.Success(c, http.StatusOK, "Users retrieved successfully", users)
}

func (h *Handler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", apperror.InvalidID())
		return
	}

	user, appErr := h.service.GetUserByID(uint(id))
	if appErr != nil {
		response.Error(c, http.StatusNotFound, "User not found", appErr)
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", user)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var dto UpdateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, err)
		return
	}

	user, appErr := h.service.UpdateUser(userID, dto)
	if appErr != nil {
		response.Error(c, http.StatusBadRequest, "Failed to update profile", appErr)
		return
	}

	response.Success(c, http.StatusOK, "Profile updated successfully", user)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", apperror.InvalidID())
		return
	}

	var dto UpdateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, err)
		return
	}

	user, appErr := h.service.UpdateUser(uint(id), dto)
	if appErr != nil {
		response.Error(c, http.StatusBadRequest, "Failed to update user", appErr)
		return
	}

	response.Success(c, http.StatusOK, "User updated successfully", user)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var dto ChangePasswordDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, err)
		return
	}

	appErr := h.service.ChangePassword(userID, dto)
	if appErr != nil {
		response.Error(c, http.StatusBadRequest, "Failed to change password", appErr)
		return
	}

	response.Success(c, http.StatusOK, "Password changed successfully", nil)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", apperror.InvalidID())
		return
	}

	appErr := h.service.DeleteUser(uint(id))
	if appErr != nil {
		response.Error(c, http.StatusNotFound, "Failed to delete user", appErr)
		return
	}

	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}
