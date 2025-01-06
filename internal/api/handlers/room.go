package handlers

import (
	"fmt"
	"net/http"
	"time"

	"oxo-game-api/internal/models"
	"oxo-game-api/pkg/utils/response"
	"oxo-game-api/pkg/utils/validator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoomHandler struct {
	db *gorm.DB
}

func NewRoomHandler(db *gorm.DB) *RoomHandler {
	return &RoomHandler{db: db}
}

// GetRooms godoc
// @Summary Get all rooms
// @Description Fetches a list of all rooms
// @Tags rooms
// @Produce json
// @Success 200 {array} models.Room
// @Failure 500 {object} response.Response
// @Router /rooms [get]
func (h *RoomHandler) GetRooms(c *gin.Context) {
	var rooms []models.Room
	if err := h.db.Find(&rooms).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch rooms")
		return
	}

	if len(rooms) == 0 {
		response.Success(c, gin.H{"message": "There is no rooms"})
		return
	}
	response.Success(c, rooms)
}

// CreateRoom godoc
// @Summary Create a new room
// @Description Creates a new room with the specified details
// @Tags rooms
// @Accept json
// @Produce json
// @Param room body models.Room true "Room information"
// @Success 200 {object} response.RoomCreateResponse"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /rooms [post]
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var input validator.RoomValidation
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	room := &models.Room{
		Name:        input.Name,
		Description: input.Description,
		Status:      "Active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.db.Create(room).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create room")
		return
	}
	response.Success(c, gin.H{"id": room.ID})
}

// GetRoomByID godoc
// @Summary Get a room by ID
// @Description Fetches a room by its ID
// @Tags rooms
// @Produce json
// @Param id path int true "Room ID"
// @Success 200 {object} models.Room
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /rooms/{id} [get]
func (h *RoomHandler) GetRoomByID(c *gin.Context) {
	id, err := validator.GetParamID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	room, err := validator.FindRoomByID(h.db, id)
	if err != nil {
		if err.Error() == "room not found" {
			response.Error(c, http.StatusNotFound, "Room not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Fail to fetch room by ID")
		return
	}
	response.Success(c, room)
}

// UpdateRoomByID godoc
// @Summary Update a room by ID
// @Description Updates the details of a room by its ID
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path int true "Room ID"
// @Param room body models.Room true "Updated room information"
// @Success 200 {object} models.Room
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /rooms/{id} [put]
func (h *RoomHandler) UpdateRoomByID(c *gin.Context) {
	id, err := validator.GetParamID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	room, err := validator.FindRoomByID(h.db, id)
	if err != nil {
		if err.Error() == "room not found" {
			response.Error(c, http.StatusNotFound, "Room not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Fail to fetch room by ID")
		return
	}

	var input validator.RoomValidation
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	room.Name = input.Name
	room.Description = input.Description
	room.Status = input.Status
	room.UpdatedAt = time.Now()

	if err := h.db.Save(&room).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update room")
		return
	}
	response.Success(c, room)
}

// DeleteRoomByID godoc
// @Summary Delete a room by ID
// @Description Deletes a room by its ID
// @Tags rooms
// @Produce json
// @Param id path int true "Room ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /rooms/{id} [delete]
func (h *RoomHandler) DeleteRoomByID(c *gin.Context) {
	id, err := validator.GetParamID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	room, err := validator.FindRoomByID(h.db, id)
	if err != nil {
		if err.Error() == "room not found" {
			response.Error(c, http.StatusNotFound, "Room not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Fail to fetch room by ID")
		return
	}
	if err := h.db.Delete(&room, id).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Fail to delete room")
		return
	}
	response.Success(c, gin.H{"message": fmt.Sprintf("Room %d is deleted successfully", id)})
}
