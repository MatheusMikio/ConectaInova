package routes

import (
	"conectainova/internal/handlers"
	"conectainova/internal/services"
	"github.com/gin-gonic/gin"
)

func initSolucaoRoutes(rg *gin.RouterGroup, service services.ISolucaoService) {
	solucao := rg.Group("solucoes")
	{
		solucao.POST("", handlers.CreateSolucaoHandler(service))
		solucao.GET("", handlers.ListSolucoesHandler(service))
		solucao.PATCH(":id", handlers.UpdateSolucaoHandler(service))
	}
}
