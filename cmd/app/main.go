package main

import (
	"csv-microservice/controllers"
	repository "csv-microservice/repositories"
	"csv-microservice/routes"
	"csv-microservice/services"
	"csv-microservice/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	utils.InitLogger()

	// Connect to PostgreSQL
	// dbConnectionString := config.GetDBConnectionString()
	dbConnectionString := `host=localhost user=postgres password=Welcome@@1234 dbname=test port=5432 sslmode=disable TimeZone=Asia/Kolkata`
	// db, err := gorm.Open(postgres.Open(dbConnectionString), &gorm.Config{})
	// if err != nil {
	// 	log.Fatalf("Failed to connect to database: %v", err)
	// }

	db := services.InitializeDatabase(dbConnectionString)
	// Initialize database (apply migrations, etc.)
	services.InitDatabase(db)

	// Initialize Gin router
	router := gin.Default()

	// Initialize layers
	repo := repository.NewRepository(db)
	service := services.NewService(repo)
	controller := controllers.NewController(service)

	// Register routes
	routes.RegisterRoutes(router, controller)

	// Start server
	log.Println("Server running on http://127.0.0.1:8080")
	router.Run(":8080")
}
