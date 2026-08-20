package main

import (
	"net/http"
	"web-crawler/controllers"
	"web-crawler/initializers"
	"web-crawler/middlewares"
	"web-crawler/models"
	"web-crawler/services"

	"github.com/gin-gonic/gin"
)

func init() { 
	initializers.LoadEnvs()
	initializers.ConnectDB()
}


func main() { 
	ipRateLimiter := models.NewIpLimiter()
	router := gin.Default()

	client := &http.Client{}
	crawlService := services.NewCrawlService(client)
	crawlHandler := controllers.NewCrawlHandler(crawlService)

	router.Use(middlewares.RateLimitingMiddleware(ipRateLimiter))
	router.Use(middlewares.StartLogging)
	
	auth := router.Group("/user")
	{
		auth.POST("/signup", controllers.CreateUser)
		auth.POST("/login", controllers.LoginUser)
		auth.GET("/profile", middlewares.CheckAuth, controllers.GetUserProfile)
	}

	protected := router.Group("/api")
	protected.Use(middlewares.CheckAuth)
	{
		auth.GET("/profile", controllers.GetUserProfile)
	}

	webcrawl := protected.Group("/crawl")
	{
		webcrawl.POST("/", crawlHandler.ProcessUrls)
	}



	router.Run()
}
