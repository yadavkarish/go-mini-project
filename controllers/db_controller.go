package controllers

import (
	"csv-microservice/services"

	"github.com/gin-gonic/gin"
)

// Controller layer to handle the request and invoke the respective service methods.
type Controller struct {
	Service services.ServiceInterface
}

func NewController(service services.ServiceInterface) *Controller {
	return &Controller{Service: service}
}

func (c *Controller) UploadCSV(ctx *gin.Context) {
	c.Service.UploadCSV(ctx)
}

func (c *Controller) ListRecords(ctx *gin.Context) {
	c.Service.ListAllEntries(ctx)
}

func (c *Controller) ListRecordsByPages(ctx *gin.Context) {
	c.Service.ListEntriesByPages(ctx)
}

func (c *Controller) SearchRecords(ctx *gin.Context) {
	c.Service.QueryUpdates(ctx)
}

func (c *Controller) AddEntries(ctx *gin.Context) {
	c.Service.AddEntries(ctx)
}

func (c *Controller) DeleteRecord(ctx *gin.Context) {
	c.Service.DeleteUpdate(ctx)
}

func (c *Controller) GetLogs(ctx *gin.Context) {
	c.Service.GetLogs(ctx)
}
