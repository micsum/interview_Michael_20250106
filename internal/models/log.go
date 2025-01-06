package models

import (
	"time"
)

type GameLog struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PlayerID  uint      `json:"player_id" gorm:"not null"`
	Action    string    `json:"action" gorm:"not null"`
	Timestamp time.Time `json:"timestamp" gorm:"not null"`
	Details   string    `json:"details" gorm:"type:text"`
}
