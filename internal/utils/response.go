package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Details string      `json:"details,omitempty"`
}

type PaginatedData struct {
	Data   interface{} `json:"data"`
	Offset int         `json:"offset"`
	Limit  int         `json:"limit"`
	Total  int         `json:"total"`
}

func SuccessResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(code, Response{
		Success: true,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, code int, message string, details ...string) {
	resp := Response{
		Success: false,
		Error:   message,
	}
	if len(details) > 0 {
		resp.Details = details[0]
	}
	c.JSON(code, resp)
}

func ValidationErrorResponse(c *gin.Context, errors []string) {
	c.JSON(400, gin.H{
		"success": false,
		"errors":  errors,
	})
}

func PaginatedResponse(c *gin.Context, code int, data interface{}, offset, limit, total int) {
	c.JSON(code, Response{
		Success: true,
		Data: PaginatedData{
			Data:   data,
			Offset: offset,
			Limit:  limit,
			Total:  total,
		},
	})
}
