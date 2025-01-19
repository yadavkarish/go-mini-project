package routes

import (
	"csv-microservice/controllers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockService satisfies services.ServiceInterface
type MockService struct{}

func (m *MockService) UploadCSV(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "UploadCSV"}) }
func (m *MockService) ListAllEntries(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "ListRecords"})
}
func (m *MockService) ListEntriesByPages(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "ListRecordsByPages"})
}
func (m *MockService) QueryUpdates(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "SearchRecords"})
}
func (m *MockService) AddRecord(ctx *gin.Context)    { ctx.JSON(200, gin.H{"message": "AddRecord"}) }
func (m *MockService) DeleteRecord(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "DeleteRecord"}) }

func TestRegisterRoutes(t *testing.T) {
	// Initialize a Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Use MockService to create a Controller
	mockService := &MockService{}
	controller := controllers.NewController(mockService)

	// Register routes
	RegisterRoutes(router, controller)

	// Define test cases
	tests := []struct {
		method   string
		path     string
		expected string
	}{
		{"POST", "/upload", "UploadCSV"},
		{"GET", "/list", "ListRecords"},
		{"GET", "/listByPages", "ListRecordsByPages"},
		{"GET", "/search", "SearchRecords"},
		{"POST", "/add", "AddRecord"},
		{"DELETE", "/delete/1", "DeleteRecord"},
	}

	// Test each route
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")
			assert.Contains(t, w.Body.String(), test.expected, "Response should contain the correct handler message")
		})
	}

	// Test for an unregistered route
	t.Run("UnregisteredRoute", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code, "Status code should be 404 for an unregistered route")
	})
}
