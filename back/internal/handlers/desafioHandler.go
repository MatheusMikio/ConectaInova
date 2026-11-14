package handlers

import (
	"net/http"

	"conectainova/internal/domain/dto/desafio"
	"conectainova/internal/domain/models"
	"conectainova/internal/response"
	"conectainova/internal/services"
	"github.com/gin-gonic/gin"
)

func CreateDesafioHandler(service services.IDesafioService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &desafio.CreateDesafioRequest{}
		if err := bindJSONObject(ctx, request); err != nil {
			response.SendError(ctx, http.StatusBadRequest, models.NewErrorMessage("JSON", "JSON inválido."))
			return
		}
		result, err := service.CreateDesafio(request)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusCreated, result)
	}
}

func ListDesafiosHandler(service services.IDesafioService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pagination, ok := getPagination(ctx)
		if !ok {
			return
		}
		result, err := service.GetDesafios(pagination.Page, pagination.Size)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusOK, result)
	}
}

func ListDesafiosByNecessidadeHandler(service services.IDesafioService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, ok := getID(ctx)
		if !ok {
			return
		}
		pagination, ok := getPagination(ctx)
		if !ok {
			return
		}
		result, err := service.GetDesafiosByNecessidade(id, pagination.Page, pagination.Size)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusOK, result)
	}
}

func UpdateDesafioHandler(service services.IDesafioService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, ok := getID(ctx)
		if !ok {
			return
		}
		request := &desafio.UpdateDesafioRequest{}
		if err := bindJSONObject(ctx, request); err != nil {
			response.SendError(ctx, http.StatusBadRequest, models.NewErrorMessage("JSON", "JSON inválido."))
			return
		}
		result, err := service.UpdateDesafio(id, request)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusOK, result)
	}
}
