package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Message string `json:"message"`
}

// OK sends a success response
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Success: true, Data: data})
}

// Failed sends an error response
func Fail(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Success: false,
		Error:   &ErrorInfo{Message: message}})
}
