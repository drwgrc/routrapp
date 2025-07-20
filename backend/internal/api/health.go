package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new health check handler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// Check handles health check requests
func (h *HealthHandler) Check(c *gin.Context) {
	health := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// Check database connection
	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil {
			health["database"] = "error"
			health["database_error"] = err.Error()
			c.JSON(http.StatusServiceUnavailable, health)
			return
		}

		if err := sqlDB.Ping(); err != nil {
			health["database"] = "error"
			health["database_error"] = err.Error()
			c.JSON(http.StatusServiceUnavailable, health)
			return
		}

		health["database"] = "ok"
	}

	c.JSON(http.StatusOK, health)
} 