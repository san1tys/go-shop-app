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

// createOrder godoc
//
// @Summary Create order for current user
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body CreateOrderInput true "Create order input"
// @Success 201 {object} Order
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/orders [post]
func (h *Handler) createOrder(c *gin.Context) {
	// userID берём из контекста, куда его положил AuthMiddleware
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

// getByID godoc
//
// @Summary Get order by ID
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} Order
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/orders/{id} [get]
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

// listMy godoc
//
// @Summary List orders for current user
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Order
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/orders/me [get]
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

	ordersList, err := h.service.ListByUser(c.Request.Context(), userID)
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

// cancel godoc
//
// @Summary Cancel order
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/orders/{id}/cancel [post]
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
