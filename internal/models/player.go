package models

import (
	"time"

	"gorm.io/gorm"
)

type Player struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:32;not null;unique"`
	LevelID   uint           `json:"level_id" gorm:"not null;default:1"`
	Level     *Level         `json:"level" gorm:"foreignKey:LevelID"`
	Balance   float64        `json:"balance" gorm:"not null;default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type Level struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"size:32;not null;unique"`
	Description string         `json:"description" gorm:"size:255"`
	MinExp      uint           `json:"min_exp" gorm:"not null;default:0"`
	MaxExp      uint           `json:"max_exp" gorm:"not null;default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (Player) TableName() string {
	return "players"
}

func (Level) TableName() string {
	return "levels"
}

func (p *Player) BeforeCreate(tx *gorm.DB) (err error) {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return
}

func (p *Player) AfterCreate(tx *gorm.DB) (err error) {
	p.UpdatedAt = time.Now()
	return
}
