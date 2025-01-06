package validator

import (
	"fmt"
	"net/http"
	"oxo-game-api/internal/models"
	"oxo-game-api/pkg/utils/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

const (
	PlayerNameRule = "required,min=3,max=100"
	LevelRule      = "required,min=1,max=100"
	RoomNameRule   = "required,min=3,max=50"
)

type CustomValidator struct {
	validator *validator.Validate
}

type Validator struct {
	db *gorm.DB
}

type PlayerValidation struct {
	Name    string `json:"name" binding:"required"`
	LevelID int    `json:"level_id"`
}

type LevelValidation struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	MinExp      int    `json:"min_exp" binding:"required"`
	MaxExp      int    `json:"max_exp" binding:"required"`
}

type RoomValidation struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Status      string `json:"status" binding:"required"`
}

func NewValidator(db *gorm.DB) *Validator {
	return &Validator{db: db}
}

func RegisterValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("player_name", validatePlayerName)
		v.RegisterValidation("level", validateLevel)
		v.RegisterValidation("room_name", validateRoomName)
	}
}

func validatePlayerName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	return len(name) >= 1 && len(name) <= 32
}

func validateRoomName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	return len(name) >= 3 && len(name) <= 50
}

func validateLevel(fl validator.FieldLevel) bool {
	level := fl.Field().Uint()
	return level >= 0 && level <= 100
}

// get and validate player/room ID
func GetParamID(c *gin.Context) (uint64, error) {
	idStr := c.Param("id")
	if idStr == "" {
		return 0, fmt.Errorf("missing ID")
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ID")
	}

	return id, nil
}

// Find player by ID
func FindPlayerByID(db *gorm.DB, id uint64) (*models.Player, error) {
	var player models.Player
	result := db.Debug().Table("players").
		Where("id = ?", id).First(&player)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("player not found")
		}
		return nil, fmt.Errorf("database error")
	}

	return &player, nil
}

// Find room by ID
func FindRoomByID(db *gorm.DB, id uint64) (*models.Room, error) {
	var room models.Room
	result := db.Debug().Table("rooms").
		Where("id = ?", id).First(&room)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("room not found")
		}
		return nil, fmt.Errorf("database error")
	}

	return &room, nil
}

// check query parameters
func CheckQueryParam(c *gin.Context, allowedParams map[string]bool) {
	for key := range c.Request.URL.Query() {
		if !allowedParams[key] {
			response.Error(c, http.StatusBadRequest, fmt.Sprintf("Unexpected query paramter: %s", key))
			c.Abort()
			return
		}
	}
}

func CheckActionValue(c *gin.Context) {
	allowedActions := map[string]bool{
		"註冊":   true,
		"登入":   true,
		"登出":   true,
		"進入房間": true,
		"退出房間": true,
		"參加挑戰": true,
		"挑戰結果": true,
	}

	if action := c.Query("action"); action != "" {
		if !allowedActions[action] {
			response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid action: %s", action))
			c.Abort()
			return
		}
	}
}
