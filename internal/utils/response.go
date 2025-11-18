package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedData struct {
	Data  interface{} `json:"data"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Total int         `json:"total"`
}

func SuccessResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(code, Response{
		Success: true,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Success: false,
		Error:   message,
	})
}

func ValidationErrorResponse(c *gin.Context, errors []string) {
	c.JSON(400, gin.H{
		"success": false,
		"errors":  errors,
	})
}

func PaginatedResponse(c *gin.Context, code int, data interface{}, page, limit, total int) {
	c.JSON(code, Response{
		Success: true,
		Data: PaginatedData{
			Data:  data,
			Page:  page,
			Limit: limit,
			Total: total,
		},
	})
}
