package handlers

import (
	"net/http"
	"oxo-game-api/internal/models"
	"oxo-game-api/pkg/utils/response"
	"oxo-game-api/pkg/utils/validator"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LevelHandler struct {
	db *gorm.DB
}

func NewLevelHandler(db *gorm.DB) *LevelHandler {
	return &LevelHandler{db: db}
}

// GetLevels godoc
// @Summary Get
// @Description Get all levels with details
// @Tags levels
// @Accept json
// @Produce json
// @Param level body models.Level true "Level information"
// @Success 200 {object} models.Level
// @Failure 500 {object} response.Response
// @Router /levels [get]
func (h *LevelHandler) GetLevels(c *gin.Context) {
	var levels []models.Level
	if err := h.db.Find(&levels).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Fail to fetch levels")
		return
	}

	response.Success(c, levels)
}

// CreateLevel godoc
// @Summary Create a new level
// @Description Create a new level with the specified name
// @Tags levels
// @Accept json
// @Produce json
// @Param level body models.Level true "Level information"
// @Success 200 {object} response.LevelCreateResponse
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /levels [post]
func (h *LevelHandler) CreateLevel(c *gin.Context) {
	var input validator.LevelValidation
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid JSON Format/missing required field")
		return
	}

	level := models.Level{
		Name:        input.Name,
		Description: input.Description,
		MinExp:      uint(input.MinExp),
		MaxExp:      uint(input.MaxExp),
		CreatedAt:   time.Now(),
	}

	if err := h.db.Create(&level).Error; err != nil {

		if strings.Contains(err.Error(), "unique constraint") {
			response.Error(c, http.StatusConflict, "Level name already exists")
			return
		}

		response.Error(c, http.StatusInternalServerError, "Fail to create level")
		return
	}

	response.Success(c, gin.H{
		"level_id": level.ID,
	})
}
