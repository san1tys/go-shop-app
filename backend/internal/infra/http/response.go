package http

import "github.com/gin-gonic/gin"

type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type APIResponse struct {
	Data any `json:"data,omitempty"`
}

func Success(c *gin.Context, data any) {
	c.JSON(200, APIResponse{Data: data})
}

func Created(c *gin.Context, data any) {
	c.JSON(201, APIResponse{Data: data})
}

func NoContent(c *gin.Context) {
	c.Status(204)
}

func Fail(c *gin.Context, status int, code, message string) {
	c.JSON(status, APIError{
		Error:   code,
		Message: message,
	})
}
