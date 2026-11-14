package routes

import (
	"conectainova/internal/handlers"
	"conectainova/internal/services"
	"github.com/gin-gonic/gin"
)

func initAvaliacaoRoutes(rg *gin.RouterGroup, service services.IAvaliacaoService) {
	avaliacao := rg.Group("avaliacoes")
	{
		avaliacao.POST("", handlers.CreateAvaliacaoHandler(service))
		avaliacao.GET("", handlers.ListAvaliacoesHandler(service))
		avaliacao.PATCH(":id", handlers.UpdateAvaliacaoHandler(service))
	}
	rg.GET("plataforma", handlers.GetPlataformaHandler())
}
