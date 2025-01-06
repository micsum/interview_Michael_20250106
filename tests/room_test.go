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

func TestCreateRoom(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDB()
	router := SetupTestRouter(db)

	t.Run("successful room creation", func(t *testing.T) {
		room := models.Room{Name: "Test Room", Description: "A room for testing"}
		body, _ := json.Marshal(room)

		req, _ := http.NewRequest(http.MethodPost, "/rooms", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var createdRoom models.Room
		err := json.Unmarshal(w.Body.Bytes(), &createdRoom)
		assert.NoError(t, err)
		assert.Equal(t, room.Name, createdRoom.Name)
		assert.Equal(t, room.Description, createdRoom.Description)
	})
}
