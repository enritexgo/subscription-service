package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"subscription-service/internal/config"
	"subscription-service/internal/handler"
	"subscription-service/internal/repository"
	"subscription-service/internal/service"
	"subscription-service/pkg/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации: ", err)
	}

	dbConn, err := db.NewPostgresDB()
	if err != nil {
		log.Fatal("Не удалось подключиться к БД: ", err)
	}
	defer dbConn.Close()

	repo := repository.NewSubscriptionRepository(dbConn)
	srv := service.NewSubscriptionService(repo)
	h := handler.NewSubscriptionHandler(srv)

	r := gin.Default()

	// Логирование
	r.Use(gin.LoggerWithWriter(log.Writer()))
	r.Use(gin.Recovery())

	// Маршруты
	api := r.Group("/api/v1")
	{
		api.POST("/subscriptions", h.Create)
		api.GET("/subscriptions/:id", h.GetByID)
		api.PUT("/subscriptions/:id", h.Update)
		api.DELETE("/subscriptions/:id", h.Delete)
		api.GET("/subscriptions", h.List)
		api.GET("/calculate", h.CalculateTotalCost)
	}

	// Swagger
	r.StaticFile("/swagger.yaml", "./swagger.yaml")
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://editor.swagger.io/?url=http://localhost:"+cfg.ServerPort+"/swagger.yaml")
	})

	log.Printf("Сервис запущен на порту :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Ошибка запуска HTTP-сервера: ", err)
	}
}
