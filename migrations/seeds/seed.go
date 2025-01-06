package seeds

import (
	"log"
	"oxo-game-api/internal/models"

	"gorm.io/gorm"
)

func SeedLevels(db *gorm.DB) error {
	defaultLevels := []models.Level{
		{
			ID:          1,
			Name:        "Beginner",
			Description: "Starting level for new players",
			MinExp:      0,
			MaxExp:      100,
		},
		{
			ID:          2,
			Name:        "Intermediate",
			Description: "Middle level",
			MinExp:      101,
			MaxExp:      200,
		},
	}

	for _, level := range defaultLevels {
		result := db.Where("name = ?", level.Name).FirstOrCreate(&level)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func SeedPlayers(db *gorm.DB) error {
	defaultPlayers := []models.Player{

		{Name: "Player 1", Balance: 0},
		{Name: "Player 2", Balance: 0},
	}
	log.Println("Creating test players:", defaultPlayers)

	for _, player := range defaultPlayers {
		result := db.Create(&player)
		if result.Error != nil {
			log.Println("Error creating player:", result.Error)
		}
	}

	return nil
}

func SeedRooms(db *gorm.DB) error {
	defaultrooms := []models.Room{
		{
			Name:        "Beginner",
			Description: "Starting level for new players",
			Status:      "Active",
		},
		{
			Name:        "Intermediate",
			Description: "Middle level",
			Status:      "Active",
		},
	}

	for _, rooms := range defaultrooms {
		result := db.Where("name = ?", rooms.Name).FirstOrCreate(&rooms)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
