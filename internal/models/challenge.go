package models

import (
	"time"
)

type Challenge struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PlayerID  uint      `json:"player_id" gorm:"not null"`
	Amount    float64   `json:"amount" gorm:"not null;default:20.01"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChallengeResult struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ChallengeID uint       `json:"challenge_id" gorm:"not null"`
	Challenge   *Challenge `json:"challenge" gorm:"foreignKey:ChallengeID"`
	PlayerID    uint       `json:"player_id" gorm:"not null"`
	Player      *Player    `json:"player" gorm:"foreignKey:PlayerID"`
	Won         bool       `json:"won" gorm:"not null"`
	CreatedAt   time.Time  `json:"created_at"`
}