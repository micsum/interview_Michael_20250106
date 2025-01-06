package models

import (
	"time"

	"gorm.io/gorm"
)

type Room struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"size:32;not null;unique"`
	Description string         `json:"description" gorm:"size:255"`
	Status      string         `json:"status" gorm:"size:20;not null;default:'available'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type Reservation struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	RoomID    uint      `json:"room_id" gorm:"not null"`
	Room      *Room     `json:"room" gorm:"foreignKey:RoomID"`
	Date      string    `json:"date" gorm:"not null"`
	Time      string    `json:"time" gorm:"not null"`
	PlayerID  uint      `json:"player_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
