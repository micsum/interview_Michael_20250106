package handlers

import (
	"net/http"
	"strconv"
	"time"

	"oxo-game-api/internal/models"
	"oxo-game-api/pkg/utils/response"
	"oxo-game-api/pkg/utils/validator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LogHandler struct {
	db *gorm.DB
}

func NewLogHandler(db *gorm.DB) *LogHandler {
	return &LogHandler{db: db}
}

// GetLogs godoc
// @Summary Get game logs
// @Description Fetches game logs based on query parameters
// @Tags logs
// @Produce json
// @Param player_id query int false "Player ID"
// @Param action query string false "Action type"
// @Param start_time query string false "Start time in RFC3339 format"
// @Param end_time query string false "End time in RFC3339 format"
// @Param limit query int false "Limit the number of logs"
// @Success 200 {array} models.GameLog
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /logs [get]
func (h *LogHandler) GetLogs(c *gin.Context) {
	allowedParams := map[string]bool{
		"limit":      true,
		"player_id":  true,
		"action":     true,
		"start_time": true,
		"end_time":   true,
	}

	validator.CheckQueryParam(c, allowedParams)
	validator.CheckActionValue(c)

	if c.IsAborted() {
		return
	}

	var logs []models.GameLog
	playerID := c.Query("player_id")
	action := c.Query("action")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	limit := c.Query("limit")

	query := h.db.Model(&models.GameLog{})

	if playerID != "" {
		id, err := strconv.ParseUint(playerID, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid player ID")
			return
		}

		if _, err := validator.FindPlayerByID(h.db, id); err != nil {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		query = query.Where("player_id = ?", playerID)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if startTime != "" {
		start, err := time.Parse(time.RFC3339, startTime)
		if err == nil {
			query = query.Where("timestamp >= ?", start)
		}
	}
	if endTime != "" {
		end, err := time.Parse(time.RFC3339, endTime)
		if err == nil {
			query = query.Where("timestamp <= ?", end)
		}
	}
	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err == nil && limitInt > 0 {
			query = query.Limit(limitInt)
		}
	}

	if err := query.Find(&logs).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Fail to fetch game logs")
		return
	}

	if len(logs) == 0 {
		response.Success(c, gin.H{"message": "There is no log of this search ing criteria"})
		return
	}

	response.Success(c, logs)
}

// CreateLog godoc
// @Summary Create a new game log
// @Description Creates a new game log entry
// @Tags logs
// @Accept json
// @Produce json
// @Param log body models.GameLog true "Game log information"
// @Success 200 {object} response.LogCreateResponse
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /logs [post]
func (h *LogHandler) CreateLog(c *gin.Context) {
	var log models.GameLog
	if err := c.ShouldBindJSON(&log); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	_, err := validator.FindPlayerByID(h.db, uint64(log.PlayerID))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Player not found")
		return
	}
	validator.CheckActionValue(c)
	if c.IsAborted() {
		return
	}
	log.Timestamp = time.Now()

	if err := h.db.Create(&log).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Fail to create game log")
		return
	}

	response.Success(c, gin.H{"log_id": log.ID})
}
