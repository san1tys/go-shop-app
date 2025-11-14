package http

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go-shop-app-backend/internal/domain"
	"go-shop-app-backend/internal/infra/auth"
	"go-shop-app-backend/internal/infra/config"
)

type Handler struct {
	db     *sql.DB
	cfg    *config.Config
	jwtMgr *auth.JWTManager
}

func NewHandler(db *sql.DB, cfg *config.Config) *Handler {
	return &Handler{
		db:     db,
		cfg:    cfg,
		jwtMgr: auth.NewJWTManager(cfg),
	}
}

// ==== AUTH ====

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type userResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// вставляем пользователя
	var u domain.User
	query := `
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, 'user')
		RETURNING id, email, role, created_at
	`
	err = h.db.QueryRow(query, req.Email, hash).Scan(&u.ID, &u.Email, &u.Role, &u.CreatedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, userResponse{
		ID:        u.ID,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	})
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	var u domain.User
	query := `
		SELECT id, email, password_hash, role, created_at
		FROM users
		WHERE email = $1
	`
	err := h.db.QueryRow(query, req.Email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query user: " + err.Error()})
		return
	}

	if err := auth.CheckPassword(u.PasswordHash, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// генерируем JWT
	token, err := h.jwtMgr.Generate(u.ID, u.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{Token: token})
}
