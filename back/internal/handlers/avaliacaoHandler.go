package handlers

import (
	"conectainova/config"
	"conectainova/internal/domain/dto/avaliacao"
	"conectainova/internal/domain/models"
	"conectainova/internal/response"
	"conectainova/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateAvaliacaoHandler(service services.IAvaliacaoService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &avaliacao.CreateAvaliacaoRequest{}
		if err := bindJSONObject(ctx, request); err != nil {
			response.SendError(ctx, 400, models.NewErrorMessage("JSON", "JSON inválido."))
			return
		}
		result, created, err := service.CreateAvaliacao(request)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		status := http.StatusOK
		if created {
			status = http.StatusCreated
		}
		response.SendSuccess(ctx, status, result)
	}
}

func ListAvaliacoesHandler(service services.IAvaliacaoService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pagination, ok := getPagination(ctx)
		if !ok {
			return
		}
		var filter struct {
			StartupID uint `form:"startup_id"`
		}
		if err := ctx.ShouldBindQuery(&filter); err != nil {
			response.SendError(ctx, 400, models.NewErrorMessage("startup_id", "Startup inválida."))
			return
		}
		result, err := service.GetAvaliacoes(filter.StartupID, pagination.Page, pagination.Size)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusOK, result)
	}
}

func UpdateAvaliacaoHandler(service services.IAvaliacaoService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, ok := getID(ctx)
		if !ok {
			return
		}
		request := &avaliacao.UpdateAvaliacaoRequest{}
		if err := bindJSONObject(ctx, request); err != nil {
			response.SendError(ctx, 400, models.NewErrorMessage("JSON", "JSON inválido."))
			return
		}
		result, err := service.UpdateAvaliacao(id, request)
		if err != nil {
			sendServiceError(ctx, err)
			return
		}
		response.SendSuccess(ctx, http.StatusOK, result)
	}
}

func GetPlataformaHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response.SendSuccess(ctx, http.StatusOK, gin.H{"associacao": config.GetAssociacaoNome(), "regiao": config.GetAssociacaoRegiao()})
	}
}
