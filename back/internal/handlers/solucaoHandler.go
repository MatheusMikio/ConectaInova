package handlers

import (
	"net/http"

	"conectainova/internal/domain/dto/solucao"
	"conectainova/internal/domain/models"
	"conectainova/internal/response"
	"conectainova/internal/services"
	"github.com/gin-gonic/gin"
)

func CreateSolucaoHandler(service services.ISolucaoService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &solucao.CreateSolucaoRequest{}
		// IDs and timestamps are never bound to persistence entities.
		if err := bindJSONObject(ctx, request); err != nil {
			response.SendError(ctx, http.StatusBadRequest, models.NewErrorMessage("JSON", "JSON inválido."))
			return
		}
		result, err := service.CreateSolucao(request)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusCreated, result)
	}
}

func ListSolucoesHandler(service services.ISolucaoService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pagination, ok := getPagination(ctx)
		if !ok {
			return
		}
		result, err := service.GetSolucoes(pagination.Page, pagination.Size)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusOK, result)
	}
}

func UpdateSolucaoHandler(service services.ISolucaoService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, ok := getID(ctx)
		if !ok {
			return
		}
		request := &solucao.UpdateSolucaoRequest{}
		if err := bindJSONObject(ctx, request); err != nil {
			response.SendError(ctx, http.StatusBadRequest, models.NewErrorMessage("JSON", "JSON inválido."))
			return
		}
		result, err := service.UpdateSolucao(id, request)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusOK, result)
	}
}
