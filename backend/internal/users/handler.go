package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-shop-app-backend/internal/domain"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", h.register)
		authGroup.POST("/login", h.login)
	}
}

func (h *Handler) register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request_body",
			"message": err.Error(),
		})
		return
	}

	resp, err := h.service.Register(c.Request.Context(), input)
	if err != nil {
		if domain.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed_to_register",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request_body",
			"message": err.Error(),
		})
		return
	}

	resp, err := h.service.Login(c.Request.Context(), input)
	if err != nil {
		if domain.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_credentials",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed_to_login",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
