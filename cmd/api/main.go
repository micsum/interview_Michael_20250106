package main

import (
	"log"
	"oxo-game-api/config"
	"oxo-game-api/internal/api/handlers"
	"oxo-game-api/migrations"
	"oxo-game-api/migrations/seeds"
	"oxo-game-api/pkg/database"
	_"oxo-game-api/docs"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

)

type Server struct {
	db  *gorm.DB
	rdb *redis.Client
}


// @title           OXO Game API
// @version         1.0

// @contact.name   Michael Sum
// @contact.url    
// @contact.email  sumchuman@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

func main() {
	
	cfg, err := config.LoadTestConfig()
	if err != nil {
		log.Fatalf("Fail to load config: %v", err)
	}

	db, err := database.InitPostgres(cfg)
	if err != nil {
		log.Fatalf("Fail to initalize database: %v", err)
	}

	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("Fail to migrate database: %v", err)
	}

	if err := seeds.SeedLevels(db); err != nil {
		log.Fatalf("Fail to seed levels: %v", err)
	}

	if err := seeds.SeedPlayers(db); err != nil {
		log.Fatalf("Fail to seed players: %v", err)
	}

	if err := seeds.SeedRooms(db); err != nil {
		log.Fatalf("Fail to seed rooms: %v", err)
	}

	playerHandler := handlers.NewPlayerHandler(db)
	levelHandler := handlers.NewLevelHandler(db)
	roomHandler := handlers.NewRoomHandler(db)
	reservationHandler := handlers.NewReservationHandler(db)
	challengeHandler := handlers.NewChallengeHandler(db)
	logHandler := handlers.NewLogHandler(db)
	paymentHandler := handlers.NewPaymentHandler(db)

	r := gin.Default()

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	players := r.Group("/players")
	{
		players.GET("", playerHandler.GetPlayers)
		players.POST("", playerHandler.CreatePlayer)
		players.GET("/:id", playerHandler.GetPlayerByID)
		players.PUT("/:id", playerHandler.UpdatePlayerByID)
		players.DELETE("/:id", playerHandler.DeletePlayerByID)
	}

	levels := r.Group("/levels")
	{
		levels.GET("", levelHandler.GetLevels)
		levels.POST("", levelHandler.CreateLevel)
	}

	rooms := r.Group("/rooms")
	{
		rooms.GET("", roomHandler.GetRooms)
		rooms.POST("", roomHandler.CreateRoom)
		rooms.GET("/:id", roomHandler.GetRoomByID)
		rooms.PUT("/:id", roomHandler.UpdateRoomByID)
		rooms.DELETE("/:id", roomHandler.DeleteRoomByID)
	}

	reservations := r.Group("/reservations")
	{
		reservations.GET("", reservationHandler.GetReservations)
		reservations.POST("", reservationHandler.CreateReservation)
	}

	challenges := r.Group("/challenges")
	{
		challenges.GET("/results", challengeHandler.GetChallengeResults)
		challenges.POST("", challengeHandler.JoinChallenge)
	}

	logs := r.Group("/logs")
	{
		logs.GET("", logHandler.GetLogs)
		logs.POST("", logHandler.CreateLog)
	}

	payments := r.Group("/payments")
	{
		payments.GET("/:id", paymentHandler.GetPayment)
		payments.POST("", paymentHandler.ProcessPayment)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Fail to run server: %v", err)
	}
}
