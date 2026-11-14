package handlers

import (
	"net/http"

	"conectainova/internal/domain/dto/lojista"
	"conectainova/internal/domain/models"
	"conectainova/internal/response"
	"conectainova/internal/services"
	"github.com/gin-gonic/gin"
)

func CreateLojistaHandler(service services.ILojistaService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &lojista.CreateLojistaRequest{}
		if err := bindJSONObject(ctx, request); err != nil {
			response.SendError(ctx, http.StatusBadRequest, models.NewErrorMessage("JSON", "JSON inválido."))
			return
		}
		result, err := service.CreateLojista(request)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusCreated, result)
	}
}

func GetLojistaHandler(service services.ILojistaService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, ok := getID(ctx)
		if !ok {
			return
		}
		result, err := service.GetLojista(id)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusOK, result)
	}
}

func UpdateLojistaHandler(service services.ILojistaService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, ok := getID(ctx)
		if !ok {
			return
		}
		request := &lojista.UpdateLojistaRequest{}
		if err := bindJSONObject(ctx, request); err != nil {
			response.SendError(ctx, http.StatusBadRequest, models.NewErrorMessage("JSON", "JSON inválido."))
			return
		}
		result, err := service.UpdateLojista(id, request)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusOK, result)
	}
}
