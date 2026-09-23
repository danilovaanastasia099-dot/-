package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"server-cooling/internal/app/handler"
	"server-cooling/internal/app/repository"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("Ошибка инициализации репозитория")
	}

	h := handler.NewHandler(repo)

	r := gin.Default()

	// Загрузка шаблонов
	r.LoadHTMLGlob("templates/*")

	// Раздача статики
	r.Static("/static", "./resources")

	// Маршруты
	r.GET("/", h.GetFeed)
	r.GET("/feed", h.GetFeed)
	r.GET("/feed/:id", h.GetFeed)
	r.GET("/add", h.GetAddPage)
	r.GET("/cards", h.GetCards)

	r.Run(":8080")
	log.Println("Server down")
}
