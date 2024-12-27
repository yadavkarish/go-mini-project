package services

import (
	"csv-microservice/models"
	repository "csv-microservice/repositories"
	"encoding/csv"
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
	AddEntries(ctx *gin.Context)
	DeleteUpdate(ctx *gin.Context)
	GetLogs(ctx *gin.Context)
}

// Implement ServiceInterface
type Service struct {
	Repo repository.RepositoryInterface
}

var db *gorm.DB
var log = logrus.New()

// Initialize PostgreSQL DB connection (Using GORM)
func InitDatabase(database *gorm.DB) {
	db = database
	db.AutoMigrate(&models.CSV{})
}

func NewService(repo repository.RepositoryInterface) *Service {
	return &Service{Repo: repo}
}

// CSV Upload and Parsing using Goroutines
func processRecords(recordChan <-chan []string, batchSize int, s *Service, wg *sync.WaitGroup) {
	defer wg.Done()
	var batch []models.CSV

	for record := range recordChan {
		recordData := models.CSV{
			SiteID:                parseInt(record[0]),
			FxiletID:              parseInt(record[1]),
			Name:                  record[2],
			Criticality:           record[3],
			RelevantComputerCount: parseInt(record[4]),
		}

		batch = append(batch, recordData)

		// Insert batch when size limit is reached
		if len(batch) >= batchSize {
			if err := s.Repo.BulkInsert(batch); err != nil {
				log.Error("Error during batch insertion: ", err)
			}
			batch = batch[:0] // Clear the batch
		}
	}

	// Insert remaining records
	if len(batch) > 0 {
		if err := s.Repo.BulkInsert(batch); err != nil {
			log.Error("Error during final batch insertion: ", err)
		}
	}
}

func (s *Service) UploadCSV(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}
	defer file.Close()

	if filepath.Ext(header.Filename) != ".csv" {
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
			log.Error("Error reading CSV row: ", err)
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

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "File uploaded and records stored"})
}

func (s *Service) ListAllEntries(ctx *gin.Context) {
	// Implementation

}

func (s *Service) ListEntriesByPages(ctx *gin.Context) {
	// Implementation
}

func (s *Service) QueryUpdates(ctx *gin.Context) {
	// Implementation
}

func (s *Service) AddEntries(ctx *gin.Context) {
	// Implementation
}

func (s *Service) DeleteUpdate(ctx *gin.Context) {
	// Implementation
}

func (s *Service) GetLogs(ctx *gin.Context) {
	// Implementation
}

// Helper function to parse integers safely
func parseInt(str string) int {
	val, _ := strconv.Atoi(strings.TrimSpace(str))
	return val
}
