package services

import (
	"csv-microservice/models"
	repository "csv-microservice/repositories"
	"csv-microservice/utils"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Service interface to define the services.
type ServiceInterface interface {
	UploadCSV(ctx *gin.Context)
	ListAllEntries(ctx *gin.Context)
	ListEntriesByPages(ctx *gin.Context)
	QueryUpdates(ctx *gin.Context)
	AddRecord(ctx *gin.Context)
	DeleteRecord(ctx *gin.Context)
	// GetLogs(ctx *gin.Context)
}

// Implement ServiceInterface
type Service struct {
	Repo repository.RepositoryInterface
}

var db *gorm.DB
var logs = logrus.New()

// Initialize PostgreSQL DB connection (Using GORM)
func InitDatabase(database *gorm.DB) {
	db = database
	db.AutoMigrate(&models.User{})
}

func NewService(repo repository.RepositoryInterface) *Service {
	return &Service{Repo: repo}
}

// CSV Upload and Parsing using Goroutines
func processRecords(recordChan <-chan []string, batchSize int, s *Service, wg *sync.WaitGroup) {
	defer wg.Done()
	var batch []models.User

	for record := range recordChan {
		recordData := models.User{
			Id:         parseInt(record[0]),
			FirstName:  record[1],
			LastName:   record[2],
			Email:      record[3],
			Age:        parseInt(record[4]),
			Gender:     record[5],
			Department: record[6],
			Company:    record[7],
			Salary:     parseFloat(record[8]),
			DateJoined: record[9],
			IsActive:   parseBool(record[10]),
		}

		batch = append(batch, recordData)

		// Insert batch when size limit is reached
		if len(batch) >= batchSize {
			if err := s.Repo.BulkInsert(batch); err != nil {
				logs.Error("Error during batch insertion: ", err)
			}
			batch = batch[:0] // Clear the batch
		}
	}

	// Insert remaining records
	if len(batch) > 0 {
		if err := s.Repo.BulkInsert(batch); err != nil {
			logs.Error("Error during final batch insertion: ", err)
		}
	}
}

func (s *Service) UploadCSV(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		utils.LogError("UploadCSV", "Failed to get file", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}
	defer file.Close()
	utils.LogInfo("UploadCSV", "Received file: "+header.Filename)

	if filepath.Ext(header.Filename) != ".csv" {
		utils.LogWarn("UploadCSV", "Invalid file format: "+header.Filename)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only CSV files are allowed."})
		return
	}

	csvReader := csv.NewReader(file)
	recordChan := make(chan []string, 1000)
	var wg sync.WaitGroup
	numWorkers := 10
	batchSize := 100 // Set batch size for bulk insertion

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processRecords(recordChan, batchSize, s, &wg)
	}

	// Read and send records to channel
	skipHeader := true
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// log.Error("Error reading CSV row: ", err)
			utils.LogError("UploadCSV", "Error reading CSV row", err)
			continue
		}
		if skipHeader {
			skipHeader = false
			continue
		}
		recordChan <- record
	}

	close(recordChan) // Signal workers to stop
	wg.Wait()         // Wait for all workers to finish

	utils.LogInfo("UploadCSV", "File processed successfully: "+header.Filename)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "File uploaded and records stored"})
}

func (s *Service) ListAllEntries(ctx *gin.Context) {
	// Implementation

}

func (s *Service) ListEntriesByPages(ctx *gin.Context) {
	// Query parameters for pagination
	pageStr := ctx.DefaultQuery("page", "1")    // Default to page 1
	limitStr := ctx.DefaultQuery("limit", "10") // Default limit to 10 items per page

	// Convert pagination parameters to integers
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	utils.LogInfo("ListEntriesByPages", fmt.Sprintf("Incoming request with page: %d, limit: %d", page, limit))
	// Calculate offset for the database query
	offset := (page - 1) * limit

	// Fetch paginated data from the database
	var entries []models.User
	result := db.Offset(offset).Limit(limit).Find(&entries)
	// fmt.Println(result)
	if result.Error != nil {
		utils.LogError("ListEntriesByPages", "Error fetching data from database", result.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch data from database",
		})
		return
	}

	// Fetch total count for metadata
	var totalCount int64
	db.Model(&models.User{}).Count(&totalCount)
	utils.LogInfo("ListEntriesByPages", fmt.Sprintf("Successfully fetched %d entries for page: %d with limit: %d", len(entries), page, limit))

	// Return paginated data
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   entries,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": totalCount,
		},
	})

	// Log the response sent to the client
	utils.LogInfo("ListEntriesByPages", fmt.Sprintf("Response sent for page: %d, limit: %d with total count: %d", page, limit, totalCount))

}

func (s *Service) QueryUpdates(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	if keyword == "" {
		utils.LogWarn("QueryUpdates", "Keyword is required but not provided")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Keyword is required",
		})
		return
	}

	// Parse pagination parameters
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "50") // Default to 50 items per page
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50 // Restrict max items per page to 100
	}
	offset := (page - 1) * limit

	utils.LogInfo("QueryUpdates", fmt.Sprintf("Request received with keyword: %s, page: %d, limit: %d", keyword, page, limit))

	// Build query parameters
	queryParams := map[string]interface{}{
		"first_name": keyword,
	}

	// Fetch matching records with pagination
	results, err := s.Repo.QueryRecords(ctx, queryParams, offset, limit)
	if err != nil {
		utils.LogError("Error", "Failed to fetch records from database", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch records",
		})
		return
	}

	// Fetch total count for metadata
	var total int64
	if err := db.Model(&models.User{}).Where("LOWER(first_name) LIKE ?", "%"+strings.ToLower(keyword)+"%").Count(&total).Error; err != nil {
		utils.LogError("QueryUpdates", "Failed to count records", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch records count",
		})
		return
	}

	// Log successful operation
	utils.LogInfo("QueryUpdates", fmt.Sprintf("Successfully fetched %d records for keyword: %s, page: %d, limit: %d", len(results), keyword, page, limit))

	// Return results with pagination metadata
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   results,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
	// Log response details
	utils.LogInfo("QueryUpdates", fmt.Sprintf("Response sent with total records: %d for keyword: %s", total, keyword))
}

func (s *Service) AddRecord(ctx *gin.Context) {
	var user models.User

	// Parse the JSON body into the User struct
	if err := ctx.ShouldBindJSON(&user); err != nil {

		logs.Warn("Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	utils.LogInfo("AddRecord", "Attempting to insert record")

	// Insert the record into the database
	err := s.Repo.InsertRecord(ctx, &user)
	if err != nil {
		utils.LogError("AddRecord", "Failed to add record to database", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to add record",
			"error":   err.Error(),
		})
		return
	}

	// Log successful insertion
	utils.LogInfo("AddRecord", "Record added successfully")
	// Respond with success
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Record added successfully",
		"data":    user,
	})
}

func (s *Service) DeleteRecord(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {

		logs.Warn("Invalid ID format", map[string]interface{}{
			"id":    idStr,
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid ID format",
		})
		return
	}

	utils.LogInfo("DeleteRecord", fmt.Sprintf("Attempting to delete record with ID: %d", id))

	// Call the repository function to delete the record
	err = s.Repo.DeleteRecord(ctx, id)
	if err != nil {
		if err.Error() == "record not found" {
			logs.Warn("Record not found", map[string]interface{}{
				"id": id,
			})
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Record not found",
			})
		} else {
			// Log and respond to unexpected errors
			utils.LogError("DeleteRecord", "Failed to delete record", fmt.Errorf("id: %d, error: %s", id, err.Error()))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to delete record",
			})
		}
		return
	}

	// Log successful deletion
	utils.LogInfo("DeleteRecord", "Record deleted successfully")

	// Respond with success if the deletion was successful
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Record deleted successfully",
	})
}

// func (s *Service) GetLogs(ctx *gin.Context) {
// 	// Implementation
// }

// Helper function to parse integers safely
func parseInt(str string) int {
	val, _ := strconv.Atoi(strings.TrimSpace(str))
	return val
}

// Helper function to parse a float64 from a string
func parseFloat(value string) float64 {
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0 // Default value if parsing fails
	}
	return num
}

// Helper function to parse a boolean from a string
func parseBool(value string) bool {
	// For simplicity, let's return true if the string is "true" (case insensitive)
	return value == "true" || value == "TRUE"
}
