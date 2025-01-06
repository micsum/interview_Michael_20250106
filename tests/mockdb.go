package tests

import (
	"fmt"
	"os"
	"testing"

	"oxo-game-api/internal/api/handlers"
	"oxo-game-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	godotenv.Load("../.env")
}

func SetupTestDB() *gorm.DB {
	host := os.Getenv("TEST_DB_HOST")
	user := os.Getenv("TEST_DB_USER")
	pass := os.Getenv("TEST_DB_PASS")
	dbname := os.Getenv("TEST_DB_NAME")
	port := os.Getenv("TEST_DB_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	db.AutoMigrate(
		&models.Player{},
		&models.Level{},
		&models.Challenge{},
		&models.ChallengeResult{},
		&models.Level{},
		&models.GameLog{},
		&models.Payment{},
		&models.Reservation{},
		&models.Room{})

	db.Exec("TRUNCATE TABLE players, challenges, levels, logs, payments, reservations, rooms RESTART IDENTITY CASCADE")
	return db
}

func SetupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	playerHandler := handlers.NewPlayerHandler(db)
	levelHandler := handlers.NewLevelHandler(db)
	roomHandler := handlers.NewRoomHandler(db)
	reservationHandler := handlers.NewReservationHandler(db)
	challengeHandler := handlers.NewChallengeHandler(db)
	logHandler := handlers.NewLogHandler(db)
	paymentHandler := handlers.NewPaymentHandler(db)
	challenges := router.Group("/challenges")
	{
		challenges.GET("/results", challengeHandler.GetChallengeResults)
		challenges.POST("", challengeHandler.JoinChallenge)
	}

	players := router.Group("/players")
	{
		players.GET("", playerHandler.GetPlayers)
		players.POST("", playerHandler.CreatePlayer)
		players.GET("/:id", playerHandler.GetPlayerByID)
		players.PUT("/:id", playerHandler.UpdatePlayerByID)
		players.DELETE("/:id", playerHandler.DeletePlayerByID)
	}

	payments := router.Group("/payments")
	{
		payments.GET("/:id", paymentHandler.GetPayment)
		payments.POST("", paymentHandler.ProcessPayment)
	}

	levels := router.Group("/levels")
	{
		levels.GET("", levelHandler.GetLevels)
		levels.POST("", levelHandler.CreateLevel)
	}

	rooms := router.Group("/rooms")
	{
		rooms.GET("", roomHandler.GetRooms)
		rooms.POST("", roomHandler.CreateRoom)
		rooms.GET("/:id", roomHandler.GetRoomByID)
		rooms.PUT("/:id", roomHandler.UpdateRoomByID)
		rooms.DELETE("/:id", roomHandler.DeleteRoomByID)
	}

	reservations := router.Group("/reservations")
	{
		reservations.GET("", reservationHandler.GetReservations)
		reservations.POST("", reservationHandler.CreateReservation)
	}

	logs := router.Group("/logs")
	{
		logs.GET("", logHandler.GetLogs)
		logs.POST("", logHandler.CreateLog)
	}
	return router
}

func CheckDatabaseConnection(db *gorm.DB) error {
	var count int64
	query := "SELECT COUNT(*) FROM players"
	if err := db.Raw(query).Scan(&count).Error; err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	fmt.Printf("Successfully connected to the database. Player count: %d\n", count)
	return nil
}

func TestDatabaseConnection(t *testing.T) {
	db := SetupTestDB()

	if db == nil {
		t.Fatal("Failed to establish a database connection")
	}

	if err := CheckDatabaseConnection(db); err != nil {
		t.Fatalf("Database connection test failed: %v", err)
	}
}
