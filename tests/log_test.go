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

func TestCreateLog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDB()
	router := SetupTestRouter(db)

	t.Run("successful log creation", func(t *testing.T) {
		log := models.GameLog{Details: "Test Log"}
		body, _ := json.Marshal(log)

		req, _ := http.NewRequest(http.MethodPost, "/logs", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var createdLog models.GameLog
		err := json.Unmarshal(w.Body.Bytes(), &createdLog)
		assert.NoError(t, err)
		assert.Equal(t, log.Details, createdLog.Details)
	})
}
