package orders

import (
	"errors"
	"net/http"
	"strconv"

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
	g := r.Group("/orders")

	g.POST("/", h.createOrder)
	g.GET("/:id", h.getByID)
	g.GET("/me", h.listMy)
	g.POST("/:id/cancel", h.cancel)
}

func (h *Handler) createOrder(c *gin.Context) {
	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "user is not authenticated",
		})
		return
	}

	userID, ok := userIDVal.(int64)
	if !ok || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "invalid user id in context",
		})
		return
	}

	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request_body",
			"message": err.Error(),
		})
		return
	}

	order, items, err := h.service.CreateOrder(c.Request.Context(), userID, input)
	if err != nil {
		if domain.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed_to_create_order",
			"message": err.Error(),
		})
		return
	}

	order.Items = items

	c.JSON(http.StatusCreated, order)
}

func (h *Handler) getByID(c *gin.Context) {
	id, err := parseIDParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "id must be a positive integer",
		})
		return
	}

	order, items, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "order_not_found",
				"message": "order not found",
			})
			return
		}

		if domain.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed_to_get_order",
			"message": err.Error(),
		})
		return
	}

	order.Items = items

	c.JSON(http.StatusOK, order)
}

func (h *Handler) listMy(c *gin.Context) {
	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "user is not authenticated",
		})
		return
	}

	userID, ok := userIDVal.(int64)
	if !ok || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "invalid user id in context",
		})
		return
	}

	page, limit, err := parsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_pagination",
			"message": err.Error(),
		})
		return
	}

	ordersList, err := h.service.ListByUser(c.Request.Context(), userID, page, limit)
	if err != nil {
		if domain.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed_to_list_orders",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ordersList)
}


func (h *Handler) cancel(c *gin.Context) {
	id, err := parseIDParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "id must be a positive integer",
		})
		return
	}

	if err := h.service.Cancel(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "order_not_found",
				"message": "order not found",
			})
			return
		}

		if domain.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed_to_cancel_order",
			"message": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func parseIDParam(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

func parsePagination(c *gin.Context) (int, int, error) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page := 1
	limit := 20

	if pageStr != "" {
		v, err := strconv.Atoi(pageStr)
		if err != nil || v <= 0 {
			return 0, 0, errors.New("page must be a positive integer")
		}
		page = v
	}

	if limitStr != "" {
		v, err := strconv.Atoi(limitStr)
		if err != nil || v <= 0 {
			return 0, 0, errors.New("limit must be a positive integer")
		}
		if v > 100 {
			return 0, 0, errors.New("limit must be less than or equal to 100")
		}
		limit = v
	}

	return page, limit, nil
}
