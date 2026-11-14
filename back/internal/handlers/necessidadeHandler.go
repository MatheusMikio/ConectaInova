package handlers

import (
	"net/http"

	"conectainova/internal/domain/dto/necessidade"
	"conectainova/internal/domain/models"
	"conectainova/internal/response"
	"conectainova/internal/services"
	"github.com/gin-gonic/gin"
)

func CreateNecessidadeHandler(service services.INecessidadeService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &necessidade.CreateNecessidadeRequest{}
		if err := bindJSONObject(ctx, request); err != nil {
			response.SendError(ctx, http.StatusBadRequest, models.NewErrorMessage("JSON", "JSON inválido."))
			return
		}
		result, err := service.CreateNecessidade(request)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusCreated, result)
	}
}

func GetNecessidadeHandler(service services.INecessidadeService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, ok := getID(ctx)
		if !ok {
			return
		}
		result, err := service.GetNecessidade(id)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusOK, result)
	}
}

func ListNecessidadesHandler(service services.INecessidadeService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pagination, ok := getPagination(ctx)
		if !ok {
			return
		}
		result, err := service.GetNecessidades(pagination.Page, pagination.Size)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusOK, result)
	}
}

func UpdateNecessidadeHandler(service services.INecessidadeService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, ok := getID(ctx)
		if !ok {
			return
		}
		request := &necessidade.UpdateNecessidadeRequest{}
		if err := bindJSONObject(ctx, request); err != nil {
			response.SendError(ctx, http.StatusBadRequest, models.NewErrorMessage("JSON", "JSON inválido."))
			return
		}
		result, err := service.UpdateNecessidade(id, request)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusOK, result)
	}
}

func GetMatchesHandler(service services.IMatchingService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, ok := getID(ctx)
		if !ok {
			return
		}
		result, err := service.GetMatches(id)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusOK, result)
	}
}
