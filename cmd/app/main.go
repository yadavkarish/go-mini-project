package main

import (
	"csv-microservice/config"
	"csv-microservice/controllers"
	repository "csv-microservice/repositories"
	"csv-microservice/routes"
	"csv-microservice/services"
	"csv-microservice/utils"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize logger
	utils.InitLogger()

	// Connect to PostgreSQL
	dbConnectionString := config.GetDBConnectionString()
	fmt.Println(dbConnectionString)
	// dbConnectionString := `host=localhost user=postgres password=Welcome@@1234 dbname=test port=5432 sslmode=disable TimeZone=Asia/Kolkata`

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
	log.Println("Server running on http://127.0.0.1:8081")
	router.Run(":8081")
}
