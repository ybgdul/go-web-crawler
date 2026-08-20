package controllers

import (
	"net/http"
	"web-crawler/models"
	"web-crawler/services"

	"github.com/gin-gonic/gin"
)

type CrawlHandler struct{
	service *services.CrawlService
}

func NewCrawlHandler(service *services.CrawlService) *CrawlHandler { 
	return &CrawlHandler{service: service}
}

func (h *CrawlHandler) ProcessUrls(c *gin.Context) { 
	var req models.CrawlRequest

	if err := c.ShouldBindJSON(&req); err != nil { 
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can't Process URL Links"})
		return
	}

	results, err := h.service.Crawl(req.Urls)

	if err != nil { 
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "Internal Server Error: " + err.Error()})
		return 
	}

	c.JSON(200, gin.H{
		"results": results,
	})
}