package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"oxo-game-api/internal/models" // Adjust the import path

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateReserv(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDB()
	router := SetupTestRouter(db)

	t.Run("successful reservation creation", func(t *testing.T) {
		reservation := models.Reservation{RoomID: 1, Date: "2023-01-01"}
		body, _ := json.Marshal(reservation)

		req, _ := http.NewRequest(http.MethodPost, "/reservations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Log("Response Body:", w.Body.String())

		var createdReservation models.Reservation
		err := json.Unmarshal(w.Body.Bytes(), &createdReservation)
		assert.NoError(t, err)
		assert.Equal(t, reservation.RoomID, createdReservation.RoomID)
		assert.Equal(t, reservation.Date, createdReservation.Date)
	})
}
