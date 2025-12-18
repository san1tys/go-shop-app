package products

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
	g := r.Group("/products")

	g.POST("/", h.create)
	g.GET("/", h.getAll)
	g.GET("/:id", h.getByID)
	g.PUT("/:id", h.update)
	g.DELETE("/:id", h.delete)
}

func (h *Handler) create(c *gin.Context) {
	var input CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request_body",
			"message": err.Error(),
		})
		return
	}

	product, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		if domain.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed_to_create_product",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *Handler) getAll(c *gin.Context) {
	page, limit, err := parsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_pagination",
			"message": err.Error(),
		})
		return
	}

	products, err := h.service.GetAll(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed_to_get_products",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, products)
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

	product, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "product_not_found",
				"message": "product not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed_to_get_product",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, product)
}


func (h *Handler) update(c *gin.Context) {
	id, err := parseIDParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "id must be a positive integer",
		})
		return
	}

	var input UpdateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request_body",
			"message": err.Error(),
		})
		return
	}

	product, err := h.service.Update(c.Request.Context(), id, input)
	if err != nil {
		if domain.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_error",
				"message": err.Error(),
			})
			return
		}

		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "product_not_found",
				"message": "product not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed_to_update_product",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, product)
}


func (h *Handler) delete(c *gin.Context) {
	id, err := parseIDParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "id must be a positive integer",
		})
		return
	}

	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "product_not_found",
				"message": "product not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed_to_delete_product",
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
