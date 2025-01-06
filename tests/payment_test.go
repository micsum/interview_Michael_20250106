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

func TestCreatePayment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDB()
	router := SetupTestRouter(db)

	t.Run("successful payment creation", func(t *testing.T) {
		payment := models.Payment{Amount: 100.0, Method: "Credit Card"}
		body, _ := json.Marshal(payment)

		req, _ := http.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var createdPayment models.Payment
		err := json.Unmarshal(w.Body.Bytes(), &createdPayment)
		assert.NoError(t, err)
		assert.Equal(t, payment.Amount, createdPayment.Amount)
		assert.Equal(t, payment.Method, createdPayment.Method)
	})
}
