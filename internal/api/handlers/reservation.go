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

type ReservationHandler struct {
	db *gorm.DB
}

func NewReservationHandler(db *gorm.DB) *ReservationHandler {
	return &ReservationHandler{db: db}
}

// GetReservations godoc
// @Summary Get reservations
// @Description Fetches reservations based on query parameters
// @Tags reservations
// @Produce json
// @Param room_id query string false "Room ID"
// @Param date query string false "Reservation date in YYYY-MM-DD format"
// @Param limit query int false "Limit the number of reservations"
// @Success 200 {array} models.Reservation
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /reservations [get]
func (h *ReservationHandler) GetReservations(c *gin.Context) {
	var reservations []models.Reservation
	allowedParams := map[string]bool{
		"room_id": true,
		"date":    true,
		"limit":   true,
	}

	validator.CheckQueryParam(c, allowedParams)

	if c.IsAborted() {
		return
	}

	roomID := c.Query("room_id")
	date := c.Query("date")
	limit := c.Query("limit")
	query := h.db.Preload("Room")

	if roomID != "" {
		query = query.Where("room_id = ?", roomID)
	}
	if date != "" {
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid date format. USE YYYY-MM-DD.")
			return
		}
		query = query.Where("date = ?", date)
	}
	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err == nil && limitInt > 0 {
			query = query.Limit(limitInt)
		}
	}

	if err := query.Find(&reservations).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Fail to fetch reservations")
		return
	}

	if len(reservations) == 0 {
		response.Success(c, gin.H{"message": "There is no reservation of this searching criteria"})
		return
	}

	response.Success(c, reservations)
}

// CreateReservation godoc
// @Summary Create a new reservation
// @Description Creates a new reservation with the specified details
// @Tags reservations
// @Accept json
// @Produce json
// @Param reservation body models.Reservation true "Reservation information"
// @Success 200 {object} response.ReservCreateResponse
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /reservations [post]
func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	var reservation models.Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	reservation.Date = time.Now().Format("2006-01-02")
	reservation.Time = time.Now().Format("15:04:05")
	reservation.CreatedAt = time.Now()
	reservation.UpdatedAt = time.Now()

	if err := h.db.Create(&reservation).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Fail to create reservation")
		return
	}
	response.Success(c, gin.H{"reservation_id": reservation.ID})
}
