package services

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// var db *gorm.DB
// var log = logrus.New()

// // Initialize PostgreSQL DB connection (Using GORM)
// func InitDatabase(database *gorm.DB) {
// 	db = database
// 	db.AutoMigrate(&models.CSV{})
// }

// CSV Upload and Parsing using Goroutines
// func UploadCSV(c *gin.Context) {
// 	// Get the file from the form
// 	file, header, err := c.Request.FormFile("file")
// 	if err != nil {
// 		log.Error("Failed to get file: ", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
// 		return
// 	}
// 	defer file.Close()

// 	// Ensure the file has a .csv extension
// 	if filepath.Ext(header.Filename) != ".csv" {
// 		log.Warn("Invalid file type uploaded")
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Only CSV files are allowed."})
// 		return
// 	}

// 	// Open the CSV reader
// 	csvReader := csv.NewReader(file)
// 	records, err := csvReader.ReadAll()
// 	if err != nil {
// 		log.Error("Failed to read CSV file: ", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read CSV"})
// 		return
// 	}

// 	// Concurrently parse and save records using goroutines
// 	var wg sync.WaitGroup
// 	for _, record := range records[1:] { // Skip header row
// 		wg.Add(1)
// 		go func(record []string) {
// 			defer wg.Done()
// 			// Validate and insert records
// 			recordData := models.CSV{
// 				SiteID:                parseInt(record[0]),
// 				FxiletID:              parseInt(record[1]),
// 				Name:                  record[2],
// 				Criticality:           record[3],
// 				RelevantComputerCount: parseInt(record[4]),
// 			}

// 			// Save to DB
// 			if err := s.Repository.AddRecord(recordData); err != nil {
// 				log.Error("Error saving record: ", err)
// 			} else {
// 				log.Info("Record saved: ", recordData)
// 			}
// 		}(record)
// 	}

// 	wg.Wait() // Wait for all goroutines to finish
// 	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "File uploaded and records stored"})
// }

// // Helper function to parse integers safely
// func parseInt(str string) int {
// 	val, _ := strconv.Atoi(strings.TrimSpace(str))
// 	return val
// }

func InitializeDatabase(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to configure database: %v", err)
	}

	sqlDB.SetMaxOpenConns(20) // Adjust based on system capacity
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Minute * 10)

	return db
}
