package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"server-cooling/internal/app/repository"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// GetFeed обрабатывает страницу ленты
func (h *Handler) GetFeed(ctx *gin.Context) {
	var service repository.Service
	// err переменная создаётся в нужных блоках и не нужна здесь

	idStr := ctx.Param("id")
	nextParam := ctx.Query("next")

	if idStr == "" {
		// Если ID не указан, показываем первую опубликованную услугу
		services, err := h.Repository.GetPublishedServices()
		if err != nil {
			logrus.Error(err)
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error": "Ошибка при получении услуг",
			})
			return
		}
		service = services[0]
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			logrus.Error(err)
			ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
				"error": "Неверный формат ID",
			})
			return
		}

		if nextParam == "true" {
			service, err = h.Repository.GetNextService(id)
		} else {
			service, err = h.Repository.GetServiceByID(id)
		}

		if err != nil {
			logrus.Error(err)
			ctx.HTML(http.StatusNotFound, "error.html", gin.H{
				"error": "Услуга не найдена",
			})
			return
		}
	}

	// управление развёрнутым описанием через параметр expanded
	expanded := ctx.Query("expanded")
	desc := service.Description
	shortDesc := desc
	maxRunes := 160
	if expanded != "true" {
		if utf8.RuneCountInString(desc) > maxRunes {
			r := []rune(desc)
			shortDesc = string(r[:maxRunes]) + "..."
		}
	}

	// формируем URL для переключения
	basePath := fmt.Sprintf("/feed/%d", service.ID)
	moreURL := basePath + "?expanded=true"
	lessURL := basePath
	if ctx.Query("next") == "true" {
		moreURL = moreURL + "&next=true"
		lessURL = lessURL + "?next=true"
	}

	totalTDP := h.Repository.CalculateTotalTDP()
	totalBTU := repository.ConvertToBTU(totalTDP)

	ctx.HTML(http.StatusOK, "feed.html", gin.H{
		"service":   service,
		"likes":     len(service.Likes),
		"totalTDP":  totalTDP,
		"totalBTU":  totalBTU,
		"shortDesc": shortDesc,
		"expanded":  expanded == "true",
		"moreURL":   moreURL,
		"lessURL":   lessURL,
	})
}

// GetAddPage обрабатывает страницу добавления
func (h *Handler) GetAddPage(ctx *gin.Context) {
	draft, err := h.Repository.GetDraftService()
	if err != nil {
		logrus.Warn("черновик не найден: ", err)
	}

	// Собираем временную корзину из параметров запроса add_id (поддерживает несколько add_id)
	addIDs := ctx.QueryArray("add_id")
	var cart []repository.Service
	for _, idStr := range addIDs {
		id, perr := strconv.Atoi(idStr)
		if perr != nil {
			logrus.Warnf("неверный id в add_id: %s", idStr)
			continue
		}
		svc, serr := h.Repository.GetServiceByID(id)
		if serr != nil {
			logrus.Warnf("услуга не найдена для add_id=%d: %v", id, serr)
			continue
		}
		cart = append(cart, svc)
	}

	ctx.HTML(http.StatusOK, "add.html", gin.H{
		"service": draft,
		"cart":    cart,
	})
}

// GetCards обрабатывает страницу плитки карточек
func (h *Handler) GetCards(ctx *gin.Context) {
	var services []repository.Service
	var err error

	minTDPStr := ctx.Query("min_tdp")
	
	if minTDPStr == "" {
		services, err = h.Repository.GetPublishedServices()
	} else {
		minTDP, err := strconv.Atoi(minTDPStr)
		if err != nil {
			logrus.Error(err)
			ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
				"error": "Неверный формат значения TDP",
			})
			return
		}
		services, err = h.Repository.FilterServicesByTDP(minTDP)
	}

	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка при получении услуг",
		})
		return
	}

	// Добавляем количество лайков для каждой услуги
	type ServiceWithLikes struct {
		Service repository.Service
		Likes   int
	}

	var servicesWithLikes []ServiceWithLikes
	for _, service := range services {
		servicesWithLikes = append(servicesWithLikes, ServiceWithLikes{
			Service: service,
			Likes:   len(service.Likes),
		})
	}

	ctx.HTML(http.StatusOK, "cards.html", gin.H{
		"services": servicesWithLikes,
		"query":    minTDPStr,
	})
}
