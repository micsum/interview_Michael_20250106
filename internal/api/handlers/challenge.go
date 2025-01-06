package handlers

import (
	"math/rand"
	"net/http"
	"time"

	"oxo-game-api/internal/models"
	"oxo-game-api/pkg/utils/response"
	"oxo-game-api/pkg/utils/validator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChallengeHandler struct {
	db *gorm.DB
}

func NewChallengeHandler(db *gorm.DB) *ChallengeHandler {
	return &ChallengeHandler{db: db}
}

// @BasePath    /

// JoinChallenge godoc
// @Summary Join a challenge
// @Description Player joins a challenge by POST, returns status
// @Tags challenges
// @Accept json
// @Produce json
// @Param challenge body models.Challenge true "Challenge information"
// @Success 200 {object} response.JoinResponse
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /challenges [post]
func (h *ChallengeHandler) JoinChallenge(c *gin.Context) {
	var challenge models.Challenge
	if err := c.ShouldBindJSON(&challenge); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var lastChallenge models.Challenge
	if err := h.db.Where("player_id = ?", challenge.PlayerID).Order("created_at desc").First(&lastChallenge).Error; err == nil {
		if time.Since(lastChallenge.CreatedAt) < time.Minute {
			response.Error(c, http.StatusForbidden, "You can only join a challenge once per minute.")
			return
		}
	}

	player, err := validator.FindPlayerByID(h.db, uint64(challenge.PlayerID))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	if player.Balance < 20.01 {
		//response.Error(c, http.StatusPaymentRequired, "Insufficient balance to join the challenge.")
		responseData := response.JoinResponse{
			Success:     false,
			ChallengeID: challenge.ID,
			Message:     "Fail to join the challenge due to insufficient balance.",
		}
	
		response.Success(c, responseData)
		return
	}

	player.Balance -= 20.01
	if err := h.db.Save(&player).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Fail to update player's balance")
		return
	}

	challenge.CreatedAt = time.Now()
	challenge.UpdatedAt = time.Now()

	if err := h.db.Create(&challenge).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Fail to join the challenge")
		return
	}

	go func(challengID uint, playerID uint) {
		time.Sleep(30 * time.Second)
		// 1% chance to win lgoic
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		won := rng.Intn(100) < 1

		result := models.ChallengeResult{
			ChallengeID: challenge.ID,
			PlayerID:    challenge.PlayerID,
			Won:         won,
			CreatedAt:   time.Now(),
		}

		if err := h.db.Create(&result).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to record challenge result")
			return
		}
	}(challenge.ID, challenge.PlayerID)
	
	responseData := response.JoinResponse{
        Success:     true,
        ChallengeID: challenge.ID,
        Message:     "Successfully joined the challenge.",
    }

    response.Success(c, responseData)

	//response.Success(c, gin.H{
	//	"success":      true,
	//	"challenge_id": challenge.ID,
	//	"message":      "Successfully to join the challenge.",
	//})
}

// GetChallengeResults godoc
// @Summary Get challenge results
// @Description Fetches all challenge results
// @Tags challenges
// @Produce json
// @Success 200 {array} models.ChallengeResult
// @Failure 500 {object} response.Response
// @Router /challenges/results [get]
func (h *ChallengeHandler) GetChallengeResults(c *gin.Context) {
	var results []models.ChallengeResult
	if err := h.db.Preload("Player").Find(&results).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch challenge results")
		return
	}
	if len(results) == 0 {
		response.Success(c, gin.H{"message": "There is no challenge result reocrd yet."})
		return
	}
	response.Success(c, results)
}
