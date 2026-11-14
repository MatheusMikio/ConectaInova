package handlers

import (
	"net/http"

	"conectainova/internal/domain/dto/startup"
	"conectainova/internal/domain/models"
	"conectainova/internal/response"
	"conectainova/internal/services"
	"github.com/gin-gonic/gin"
)

func CreateStartupHandler(service services.IStartupService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &startup.CreateStartupRequest{}
		if err := bindJSONObject(ctx, request); err != nil {
			response.SendError(ctx, http.StatusBadRequest, models.NewErrorMessage("JSON", "JSON inválido."))
			return
		}
		result, err := service.CreateStartup(request)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusCreated, result)
	}
}

func ListStartupsHandler(service services.IStartupService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pagination, ok := getPagination(ctx)
		if !ok {
			return
		}
		result, err := service.GetStartups(pagination.Page, pagination.Size)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusOK, result)
	}
}

func UpdateStartupHandler(service services.IStartupService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, ok := getID(ctx)
		if !ok {
			return
		}
		request := &startup.UpdateStartupRequest{}
		if err := bindJSONObject(ctx, request); err != nil {
			response.SendError(ctx, http.StatusBadRequest, models.NewErrorMessage("JSON", "JSON inválido."))
			return
		}
		result, err := service.UpdateStartup(id, request)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusOK, result)
	}
}
