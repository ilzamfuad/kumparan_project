package main

import (
	"fmt"
	"kumparan_project/config"
	"kumparan_project/controller"
	"kumparan_project/middleware"
	"kumparan_project/repository"
	"kumparan_project/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	router := gin.Default()

	db := config.BuildDB()
	defer func() {
		if err := db.Ping(); err != nil {
			panic(err)
		} else {
			_ = db.Close()
		}
	}()

	redis := config.InitRedis()
	defer func() {
		if err := redis.Ping(); err != nil {
			panic(err)
		} else {
			_ = redis.Close()
		}
	}()

	articleRepo := repository.NewArticleRepository(db)
	articleService := service.NewArticleService(articleRepo)
	articleController := controller.NewArticleController(articleService, redis)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "success",
		})
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	protected := router.Group("/articles")
	protected.Use(middleware.RateLimitMiddlewareWithRedis(redis))

	protected.GET("/", articleController.GetArticles)
	protected.POST("/", articleController.CreateArticles)

	server := fmt.Sprintf(":%s", os.Getenv("PORT"))
	router.Run(server)

}
