package response

import (
	"github.com/ardipermana59/go-template/internal/common/apperror"
	customValidator "github.com/ardipermana59/go-template/pkg/validator"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string, errors apperror.AppErrors) {
	c.JSON(code, Response{
		Success: false,
		Message: message,
		Error:   errors,
	})
}

func ValidationError(c *gin.Context, err error) {
	errors := customValidator.FormatValidationErrors(err)

	c.JSON(400, Response{
		Success: false,
		Message: "Validation failed",
		Error:   errors,
	})
}

func InternalError(c *gin.Context, err error) {
	c.JSON(500, Response{
		Success: false,
		Message: "Internal server error",
		Error:   apperror.DatabaseError(err),
	})
}
