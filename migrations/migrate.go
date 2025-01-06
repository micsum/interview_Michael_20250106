package migrations

import (
	"oxo-game-api/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Level{},
		&models.Player{},
		&models.Reservation{},
		&models.Room{},
		&models.Challenge{},
		&models.ChallengeResult{},
		&models.GameLog{},
		&models.Payment{},
	)
}
