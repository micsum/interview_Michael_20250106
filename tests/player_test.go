package tests

import (
	"encoding/json"

	"net/http"
	"net/http/httptest"

	"oxo-game-api/internal/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetPlayers(t *testing.T) {
	db := SetupTestDB()
	router := SetupTestRouter(db)

	testPlayers := []models.Player{
		{Name: "Player 1", Balance: 0},
		{Name: "Player 2", Balance: 0},
	}

	t.Log("Creating test players:", testPlayers)

	for _, player := range testPlayers {
		result := db.Create(&player)
		if result.Error != nil {
			t.Log("Error creating player:", result.Error)
		}
	}

	var dbPlayers []models.Player
	db.Find(&dbPlayers)
	t.Log("Players in database:", dbPlayers)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/players", nil)
	router.ServeHTTP(w, req)

	t.Log("Response Status:", w.Code)
	t.Log("Response Body:", w.Body.String())

	var response struct {
		Data []models.Player `json:"data"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Log("Error decoding response:", err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, err)
	assert.Equal(t, len(testPlayers), len(response.Data))
	assert.Equal(t, testPlayers[0].Name, response.Data[0].Name)
	assert.Equal(t, testPlayers[1].Name, response.Data[1].Name)
}

func TestGetPlayerByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDB()
	router := SetupTestRouter(db)

	player := models.Player{Name: "Player 8"}
	db.Create(&player)

	t.Run("successful fetch of player by ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/players/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var fetchedPlayer models.Player
		err := json.Unmarshal(w.Body.Bytes(), &fetchedPlayer)
		assert.NoError(t, err)
		assert.Equal(t, "Player One", fetchedPlayer.Name)
	})

	t.Run("player not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/players/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		t.Log("Response Body:", w.Body.String())

		var response struct {
			Code    int         `json:"code"`
			Message string      `json:"message"`
			Data    interface{} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 404, response.Code)
		assert.Equal(t, "Player not found", response.Message)
		assert.Nil(t, response.Data)
	})
}
