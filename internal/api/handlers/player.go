package handlers

import (
	"fmt"
	"log"
	"net/http"
	"oxo-game-api/internal/models"
	"oxo-game-api/pkg/utils/response"
	"oxo-game-api/pkg/utils/validator"

	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PlayerHandler struct {
	db *gorm.DB
}

func NewPlayerHandler(db *gorm.DB) *PlayerHandler {
	return &PlayerHandler{
		db: db,
	}
}

// GetPlayers godoc
// @Summary Get all players
// @Description Fetches a list of all players
// @Tags players
// @Produce json
// @Success 200 {array} models.Player
// @Failure 500 {object} response.Response
// @Router /players [get]
func (h *PlayerHandler) GetPlayers(c *gin.Context) {
	var players []models.Player
	if err := h.db.Find(&players).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Fail to fetch players")
		log.Printf("Eror fetching players: %v", err)
		return
	}

	response.Success(c, players)
}

// CreatePlayer godoc
// @Summary Create a new player
// @Description Creates a new player with the specified details
// @Tags players
// @Accept json
// @Produce json
// @Param player body models.Player true "Player information"
// @Success 200 {object} response.PlayerCreateResponse
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /players [post]
func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var input validator.PlayerValidation

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	player := &models.Player{
		Name:      input.Name,
		LevelID:   uint(input.LevelID),
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.db.Create(player).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Fail to create player")
		return
	}

	response.Success(c, gin.H{
		"id": player.ID,
	})
}

// GetPlayerByID godoc
// @Summary Get a player by ID
// @Description Fetches a player by their ID
// @Tags players
// @Produce json
// @Param id path int true "Player ID"
// @Success 200 {object} models.Player
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /players/{id} [get]
func (h *PlayerHandler) GetPlayerByID(c *gin.Context) {
	id, err := validator.GetParamID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	player, err := validator.FindPlayerByID(h.db, id)
	if err != nil {
		if err.Error() == "player not found" {
			response.Error(c, http.StatusNotFound, "Player not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Fail to fetch player by ID")
		return
	}

	response.Success(c, player)
}

// UpdatePlayerByID godoc
// @Summary Update a player by ID
// @Description Updates the details of a player by their ID
// @Tags players
// @Accept json
// @Produce json
// @Param id path int true "Player ID"
// @Param player body models.Player true "Updated player information"
// @Success 200 {object} models.Player
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /players/{id} [put]
func (h *PlayerHandler) UpdatePlayerByID(c *gin.Context) {
	id, err := validator.GetParamID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	player, err := validator.FindPlayerByID(h.db, id)
	if err != nil {
		if err.Error() == "player not found" {
			response.Error(c, http.StatusNotFound, "Player not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Fail to fetch player by ID")
		return
	}

	var input validator.PlayerValidation
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	var level models.Level
	if err := h.db.First(&level, input.LevelID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, http.StatusBadRequest, "Invalid level ID")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to verify level")
		return
	}
	player.Name = input.Name
	player.LevelID = uint(input.LevelID)

	if err := h.db.Save(player).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Fail to update player")
		return
	}

	response.Success(c, player)
}

// DeletePlayerByID godoc
// @Summary Delete a player by ID
// @Description Deletes a player by their ID
// @Tags players
// @Produce json
// @Param id path int true "Player ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /players/{id} [delete]
func (h *PlayerHandler) DeletePlayerByID(c *gin.Context) {
	id, err := validator.GetParamID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	player, err := validator.FindPlayerByID(h.db, id)
	if err != nil {
		if err.Error() == "player not found" {
			response.Error(c, http.StatusNotFound, "Player not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Fail to fetch player by ID")
		return
	}

	if err := h.db.Delete(&player).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Fail to delete player")
		return
	}

	response.Success(c, gin.H{"message": fmt.Sprintf("Player %d is deleted successfully", id)})
}
