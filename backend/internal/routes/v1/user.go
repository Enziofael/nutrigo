package apiv1

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Enziofael/nutrigo/backend/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Status string

const (
	StatusRequested  Status = "requested"
	StatusConfirmed  Status = "confirmed"
	StatusRestricted Status = "restricted"
	StatusBanned     Status = "banned"
	StatusAdmin      Status = "admin"
)

type User struct {
	ID             uint      `json:"id"`
	TgID           uint      `json:"tgID"`
	TgTag          string    `json:"tgTag"`
	CreatedAt      time.Time `json:"createdAt"`
	LastMessagedAt time.Time `json:"lastMessagedAt"`
	Status         string    `json:"status"`
}

type UserHandler struct {
	DB *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func RegisterUserRoutes(group *gin.RouterGroup, s *server.Server) {
	uh := NewUserHandler(s.DB)

	user := group.Group("/user")
	{
		user.GET("/", uh.list)
		user.GET("/:tg_id", uh.getById)
	}
	{
		user.POST("/", uh.create)
	}
}

func (h *UserHandler) list(c *gin.Context) {
	rows, err := h.DB.Query("SELECT id, tg_id, tg_tag, created_at, last_messaged_at, status FROM users")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.TgID, &u.TgTag, &u.CreatedAt, &u.LastMessagedAt, &u.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) getById(c *gin.Context) {
	var u User
	err := h.DB.QueryRowContext(c.Request.Context(), `
	SELECT id, tg_id, tg_tag, created_at, last_messaged_at, status FROM users
	WHERE tg_id = $1`, c.Param("tg_id")).
		Scan(&u.ID, &u.TgID, &u.TgTag, &u.CreatedAt, &u.LastMessagedAt, &u.Status)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) create(c *gin.Context) {
	var req struct {
		TgID  uint   `json:"tgId" binding:"required"`
		TgTag string `json:"tgTag" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	err := h.DB.QueryRowContext(c.Request.Context(), `
        INSERT INTO users (tg_id, tg_tag)
        VALUES ($1, $2)
        RETURNING id, tg_id, tg_tag, created_at, last_messaged_at, status
    `, req.TgID, req.TgTag).Scan(
		&user.ID,
		&user.TgID,
		&user.TgTag,
		&user.CreatedAt,
		&user.LastMessagedAt,
		&user.Status,
	)
	if err != nil {
		// Проверяем на уникальность
		if isUniqueViolation(err) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "user with this tg_id or tg_tag already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func isUniqueViolation(err error) bool {
	// Ошибка PostgreSQL: "23505" = unique_violation
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}
