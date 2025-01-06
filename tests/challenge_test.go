package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"oxo-game-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJoinChallenge(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDB()
	router := SetupTestRouter(db)

	t.Run("successful join", func(t *testing.T) {
		challenge := models.Challenge{PlayerID: 2}
		body, _ := json.Marshal(challenge)

		req, _ := http.NewRequest(http.MethodPost, "/challenges", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var createdChallenge models.Challenge
		if err := db.First(&createdChallenge, "player_id = ?", 2).Error; err != nil {
			t.Fatalf("could not find created challenge: %v", err)
		}
		assert.Equal(t, 2, createdChallenge.PlayerID)
	})
}

func TestGetChallengeResults(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDB()
	router := SetupTestRouter(db)

	t.Run("successful fetch", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/challenges/results", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

}
