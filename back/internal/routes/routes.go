package routes

import (
	"conectainova/internal/container"
	"github.com/gin-gonic/gin"
)

func initRoutes(router *gin.Engine, container *container.Container) {
	basePath := "api/v1"

	v1 := router.Group(basePath)
	{
		publicRoutes := v1.Group("")
		{
			initLojistaRoutes(publicRoutes, container.LojistaService)
			initNecessidadeRoutes(publicRoutes, container.NecessidadeService, container.MatchingService, container.DesafioService)
			initStartupRoutes(publicRoutes, container.StartupService)
			initSolucaoRoutes(publicRoutes, container.SolucaoService)
			initDesafioRoutes(publicRoutes, container.DesafioService)
			initAvaliacaoRoutes(publicRoutes, container.AvaliacaoService)
		}
	}
}
