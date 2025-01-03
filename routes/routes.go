package routes

import (
	"csv-microservice/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all API routes and maps them to the respective controller methods.
func RegisterRoutes(router *gin.Engine, controller *controllers.Controller) {
	router.POST("/upload", controller.UploadCSV)
	router.GET("/list", controller.ListRecords)
	router.GET("/listByPages", controller.ListRecordsByPages)
	router.GET("/search", controller.SearchRecords)
	router.POST("/add", controller.AddRecord)
	router.DELETE("/delete/:id", controller.DeleteRecord)
	router.GET("/logs", controller.GetLogs)
}
