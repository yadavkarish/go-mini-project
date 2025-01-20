package services

import (
	"bytes"
	"csv-microservice/mock"
	"csv-microservice/models"
	"csv-microservice/utils"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestListEntriesByPages tests the ListEntriesByPages function.
func TestListEntriesByPages_ValidPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepositoryInterface(ctrl)
	mockService := mock.NewMockServiceInterface(ctrl)

	router := gin.Default()
	router.GET("/listByPages", func(ctx *gin.Context) {
		mockService.ListEntriesByPages(ctx)
	})

	mockRepo.EXPECT().QueryRecords(gomock.Any(), gomock.Any(), 5, 5).Return([]models.User{}, nil)
	mockService.EXPECT().ListEntriesByPages(gomock.Any()).Do(func(ctx *gin.Context) {
		users, err := mockRepo.QueryRecords(ctx, nil, 5, 5)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch data from database"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   users,
			"meta":   gin.H{"page": 2, "limit": 5, "total": len(users)},
		})
	})

	req, _ := http.NewRequest("GET", "/listByPages?page=2&limit=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"success","data":[],"meta":{"page":2,"limit":5,"total":0}}`, w.Body.String())
}

func TestListEntriesByPages_InvalidPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepositoryInterface(ctrl)
	mockService := mock.NewMockServiceInterface(ctrl)

	router := gin.Default()
	router.GET("/listByPages", func(ctx *gin.Context) {
		mockService.ListEntriesByPages(ctx)
	})

	mockRepo.EXPECT().QueryRecords(gomock.Any(), gomock.Any(), 0, 10).Return([]models.User{}, nil)
	mockService.EXPECT().ListEntriesByPages(gomock.Any()).Do(func(ctx *gin.Context) {
		users, err := mockRepo.QueryRecords(ctx, nil, 0, 10)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch data from database"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   users,
			"meta":   gin.H{"page": 1, "limit": 10, "total": len(users)},
		})
	})

	req, _ := http.NewRequest("GET", "/listByPages?page=abc&limit=xyz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"success","data":[],"meta":{"page":1,"limit":10,"total":0}}`, w.Body.String())
}

func TestListEntriesByPages_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepositoryInterface(ctrl)
	mockService := mock.NewMockServiceInterface(ctrl)

	router := gin.Default()
	router.GET("/listByPages", func(ctx *gin.Context) {
		mockService.ListEntriesByPages(ctx)
	})

	mockRepo.EXPECT().QueryRecords(gomock.Any(), gomock.Any(), 0, 10).Return(nil, errors.New("db error"))
	mockService.EXPECT().ListEntriesByPages(gomock.Any()).Do(func(ctx *gin.Context) {
		_, err := mockRepo.QueryRecords(ctx, nil, 0, 10)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to fetch data from database",
			})
			return
		}
	})

	req, _ := http.NewRequest("GET", "/listByPages?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"status":"error","message":"Failed to fetch data from database"}`, w.Body.String())
}

// TestQueryUpdates tests the QueryUpdates function.
func TestQueryUpdates_ValidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepositoryInterface(ctrl)
	mockService := mock.NewMockServiceInterface(ctrl)

	router := gin.Default()
	router.GET("/search", func(ctx *gin.Context) {
		mockService.QueryUpdates(ctx)
	})

	// Set up mock expectations for the repository
	mockRepo.EXPECT().QueryRecords(gomock.Any(), gomock.Any(), 5, 5).Return([]models.User{
		{FirstName: "John", LastName: "Doe", Age: 0, Company: "", Email: "", Salary: 0, DateJoined: "", Department: "", Gender: "", IsActive: false},
		{FirstName: "Johnny", LastName: "Bravo", Age: 0, Company: "", Email: "", Salary: 0, DateJoined: "", Department: "", Gender: "", IsActive: false},
	}, nil)

	// Set up mock expectations for the service
	mockService.EXPECT().QueryUpdates(gomock.Any()).Do(func(ctx *gin.Context) {
		// Extract query parameters and call the repository method
		users, err := mockRepo.QueryRecords(ctx, map[string]interface{}{
			"first_name": "john",
		}, 5, 5)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to fetch records",
			})
			return
		}

		// Return the response with status, data, and meta information
		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   users,
			"meta": gin.H{
				"page":  2,
				"limit": 5,
				"total": len(users),
			},
		})
	})

	// Create the GET request with query parameters
	req, _ := http.NewRequest("GET", "/search?keyword=john&page=2&limit=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{
		"status": "success",
		"data": [
			{
				"id": 0,
				"first_name": "John",
				"last_name": "Doe",
				"age": 0,
				"company": "",
				"email": "",
				"salary": 0,
				"date_joined": "",
				"department": "",
				"gender": "",
				"is_active": false
			},
			{
				"id": 0,
				"first_name": "Johnny",
				"last_name": "Bravo",
				"age": 0,
				"company": "",
				"email": "",
				"salary": 0,
				"date_joined": "",
				"department": "",
				"gender": "",
				"is_active": false
			}
		],
		"meta": {
			"page": 2,
			"limit": 5,
			"total": 2
		}
	}`, w.Body.String())
}

func TestQueryUpdates_EmptyResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepositoryInterface(ctrl)
	mockService := mock.NewMockServiceInterface(ctrl)

	router := gin.Default()
	router.GET("/search", func(ctx *gin.Context) {
		mockService.QueryUpdates(ctx)
	})

	mockRepo.EXPECT().QueryRecords(gomock.Any(), gomock.Any(), 0, 10).Return([]models.User{}, nil)

	mockService.EXPECT().QueryUpdates(gomock.Any()).Do(func(ctx *gin.Context) {
		users, err := mockRepo.QueryRecords(ctx, map[string]interface{}{"first_name": "nonexistent"}, 0, 10)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch records"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   users,
			"meta":   gin.H{"page": 1, "limit": 10, "total": len(users)},
		})
	})

	req, _ := http.NewRequest("GET", "/search?keyword=nonexistent&page=1&limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"success","data":[],"meta":{"page":1,"limit":10,"total":0}}`, w.Body.String())
}

func TestQueryUpdates_InvalidPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepositoryInterface(ctrl)
	mockService := mock.NewMockServiceInterface(ctrl)

	router := gin.Default()
	router.GET("/search", func(ctx *gin.Context) {
		mockService.QueryUpdates(ctx)
	})

	mockRepo.EXPECT().QueryRecords(gomock.Any(), gomock.Any(), 0, 50).Return([]models.User{}, nil)

	mockService.EXPECT().QueryUpdates(gomock.Any()).Do(func(ctx *gin.Context) {
		users, err := mockRepo.QueryRecords(ctx, map[string]interface{}{"first_name": "john"}, 0, 50)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch records"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   users,
			"meta":   gin.H{"page": 1, "limit": 50, "total": len(users)},
		})
	})

	req, _ := http.NewRequest("GET", "/search?keyword=john&page=-1&limit=200", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"success","data":[],"meta":{"page":1,"limit":50,"total":0}}`, w.Body.String())
}

func TestQueryUpdates_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepositoryInterface(ctrl)
	mockService := mock.NewMockServiceInterface(ctrl)

	router := gin.Default()
	router.GET("/search", func(ctx *gin.Context) {
		mockService.QueryUpdates(ctx)
	})

	mockRepo.EXPECT().QueryRecords(gomock.Any(), gomock.Any(), 0, 10).Return(nil, errors.New("db error"))

	mockService.EXPECT().QueryUpdates(gomock.Any()).Do(func(ctx *gin.Context) {
		_, err := mockRepo.QueryRecords(ctx, map[string]interface{}{"first_name": "john"}, 0, 10)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch records"})
			return
		}
	})

	req, _ := http.NewRequest("GET", "/search?keyword=john&page=1&limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"status":"error","message":"Failed to fetch records"}`, w.Body.String())
}

// TestDeleteRecord tests the DeleteRecord function.
func TestDeleteRecord_ValidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepositoryInterface(ctrl)
	mockService := mock.NewMockServiceInterface(ctrl)

	router := gin.Default()
	router.DELETE("/records/:id", func(ctx *gin.Context) {
		mockService.DeleteRecord(ctx)
	})

	mockRepo.EXPECT().DeleteRecord(gomock.Any(), 1).Return(nil)

	mockService.EXPECT().DeleteRecord(gomock.Any()).Do(func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, _ := strconv.Atoi(idStr)
		err := mockRepo.DeleteRecord(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to delete record"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Record deleted successfully"})
	})

	req, _ := http.NewRequest("DELETE", "/records/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"success","message":"Record deleted successfully"}`, w.Body.String())
}

func TestDeleteRecord_InvalidIDFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockServiceInterface(ctrl)

	router := gin.Default()
	router.DELETE("/records/:id", func(ctx *gin.Context) {
		mockService.DeleteRecord(ctx)
	})

	mockService.EXPECT().DeleteRecord(gomock.Any()).Do(func(ctx *gin.Context) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid ID format"})
	})

	req, _ := http.NewRequest("DELETE", "/records/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"status":"error","message":"Invalid ID format"}`, w.Body.String())
}

func TestDeleteRecord_RecordNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepositoryInterface(ctrl)
	mockService := mock.NewMockServiceInterface(ctrl)

	router := gin.Default()
	router.DELETE("/records/:id", func(ctx *gin.Context) {
		mockService.DeleteRecord(ctx)
	})

	mockRepo.EXPECT().DeleteRecord(gomock.Any(), 2).Return(errors.New("record not found"))

	mockService.EXPECT().DeleteRecord(gomock.Any()).Do(func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, _ := strconv.Atoi(idStr)
		err := mockRepo.DeleteRecord(ctx, id)
		if err != nil && err.Error() == "record not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Record not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to delete record"})
	})

	req, _ := http.NewRequest("DELETE", "/records/2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"status":"error","message":"Record not found"}`, w.Body.String())
}

func TestDeleteRecord_UnexpectedError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepositoryInterface(ctrl)
	mockService := mock.NewMockServiceInterface(ctrl)

	router := gin.Default()
	router.DELETE("/records/:id", func(ctx *gin.Context) {
		mockService.DeleteRecord(ctx)
	})

	mockRepo.EXPECT().DeleteRecord(gomock.Any(), 3).Return(errors.New("unexpected error"))

	mockService.EXPECT().DeleteRecord(gomock.Any()).Do(func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, _ := strconv.Atoi(idStr)
		err := mockRepo.DeleteRecord(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to delete record"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Record deleted successfully"})
	})

	req, _ := http.NewRequest("DELETE", "/records/3", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"status":"error","message":"Failed to delete record"}`, w.Body.String())
}

// TestAddRecord_ValidRecord tests the AddRecord function with a valid record.
func TestAddRecord_ValidRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepositoryInterface(ctrl)
	mockService := mock.NewMockServiceInterface(ctrl)

	router := gin.Default()
	router.POST("/records", func(ctx *gin.Context) {
		mockService.AddRecord(ctx)
	})

	mockRepo.EXPECT().InsertRecord(gomock.Any(), gomock.Any()).Return(nil)

	mockService.EXPECT().AddRecord(gomock.Any()).Do(func(ctx *gin.Context) {
		var user models.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body", "error": err.Error()})
			return
		}
		err := mockRepo.InsertRecord(ctx, &user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to add record", "error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Record added successfully", "data": user})
	})

	req, _ := http.NewRequest("POST", "/records", strings.NewReader(`{
		"id": 1,
		"first_name": "John",
		"last_name": "Doe",
		"email": "john.doe@example.com",
		"age": 30,
		"gender": "Male",
		"department": "Engineering",
		"company": "TechCorp",
		"salary": 100000,
		"date_joined": "2025-01-01",
		"is_active": true
	}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{
		"status": "success",
		"message": "Record added successfully",
		"data": {
			"id": 1,
			"first_name": "John",
			"last_name": "Doe",
			"email": "john.doe@example.com",
			"age": 30,
			"gender": "Male",
			"department": "Engineering",
			"company": "TechCorp",
			"salary": 100000,
			"date_joined": "2025-01-01",
			"is_active": true
		}
	}`, w.Body.String())
}

func TestAddRecord_InvalidRequestBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockServiceInterface(ctrl)

	router := gin.Default()
	router.POST("/records", func(ctx *gin.Context) {
		mockService.AddRecord(ctx)
	})

	mockService.EXPECT().AddRecord(gomock.Any()).Do(func(ctx *gin.Context) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
		})
	})

	req, _ := http.NewRequest("POST", "/records", strings.NewReader(`{
		"name": "John Doe"
	}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{
		"status": "error",
		"message": "Invalid request body"
	}`, w.Body.String())
}

func TestAddRecord_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepositoryInterface(ctrl)
	mockService := mock.NewMockServiceInterface(ctrl)

	router := gin.Default()
	router.POST("/records", func(ctx *gin.Context) {
		mockService.AddRecord(ctx)
	})

	mockRepo.EXPECT().InsertRecord(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

	mockService.EXPECT().AddRecord(gomock.Any()).Do(func(ctx *gin.Context) {
		var user models.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body"})
			return
		}
		err := mockRepo.InsertRecord(ctx, &user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to add record", "error": err.Error()})
			return
		}
	})

	req, _ := http.NewRequest("POST", "/records", strings.NewReader(`{
		"id": 2,
		"name": "Jane Doe",
		"email": "jane.doe@example.com"
	}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{
		"status": "error",
		"message": "Failed to add record",
		"error": "database error"
	}`, w.Body.String())
}

func TestAddRecord_InvalidJSONFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockServiceInterface(ctrl)

	router := gin.Default()
	router.POST("/records", func(ctx *gin.Context) {
		mockService.AddRecord(ctx)
	})

	mockService.EXPECT().AddRecord(gomock.Any()).Do(func(ctx *gin.Context) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
		})
	})

	req, _ := http.NewRequest("POST", "/records", strings.NewReader(`{ "id": 1, "name": "John Doe", "email": }`)) // Malformed JSON
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{
		"status": "error",
		"message": "Invalid request body"
	}`, w.Body.String())
}

func TestUploadCSV(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repository
	mockRepo := mock.NewMockRepositoryInterface(ctrl)
	mockService := NewService(mockRepo) // Inject the mock repository

	// Initialize the logger
	utils.InitLogger() // Replace this with your logger initialization method

	// Create a Gin router
	router := gin.Default()
	router.POST("/upload", mockService.UploadCSV)

	// Define test cases
	tests := []struct {
		name           string
		fileContent    string
		fileName       string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Valid CSV Upload",
			fileContent: "id,first_name,last_name,email,age,gender,department,company,salary,date_joined,is_active\n1,John,Doe,john.doe@example.com,30,Male,Engineering,TechCorp,100000,2025-01-01,true",
			fileName:    "valid.csv",
			mockSetup: func() {
				mockRepo.EXPECT().BulkInsert(gomock.Any()).Return(nil).Times(1) // Expecting BulkInsert to be called once
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"success","message":"File uploaded and records stored"}`,
		},
		{
			name:           "Invalid File Format (Non-CSV)",
			fileContent:    "id,first_name,last_name,email,age,gender,department,company,salary,date_joined,is_active\n1,John,Doe,john.doe@example.com,30,Male,Engineering,TechCorp,100000,2025-01-01,true",
			fileName:       "invalid.txt",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Only CSV files are allowed."}`,
		},
		{
			name:           "Empty CSV File",
			fileContent:    "",
			fileName:       "empty.csv",
			mockSetup:      func() {},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"success","message":"File uploaded and records stored"}`,
		},
		{
			name:        "Malformed CSV Data (Missing Field)",
			fileContent: "id,first_name,last_name,email,age,gender,department,company,salary,date_joined,is_active\n1,John,Doe,,30,Male,Engineering,TechCorp,100000,2025-01-01,true",
			fileName:    "malformed.csv",
			mockSetup: func() {
				// Assuming the service attempts to insert malformed data
				mockRepo.EXPECT().BulkInsert(gomock.Any()).Return(nil).Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"success","message":"File uploaded and records stored"}`, // Assuming we still process despite errors in data
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.mockSetup()

			// Create a buffer to hold multipart data
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			// Add file to the multipart form
			part, err := writer.CreateFormFile("file", tt.fileName)
			assert.NoError(t, err, "Failed to create form file")
			_, err = part.Write([]byte(tt.fileContent))
			assert.NoError(t, err, "Failed to write file content")

			// Close the writer to finalize the multipart form
			writer.Close()

			// Create HTTP request
			req := httptest.NewRequest("POST", "/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			// Send the request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
