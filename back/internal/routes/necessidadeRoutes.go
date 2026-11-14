package routes

import (
	"conectainova/internal/handlers"
	"conectainova/internal/services"
	"github.com/gin-gonic/gin"
)

func initNecessidadeRoutes(rg *gin.RouterGroup, service services.INecessidadeService, matchingService services.IMatchingService, desafioService services.IDesafioService) {
	necessidade := rg.Group("necessidades")
	{
		necessidade.POST("", handlers.CreateNecessidadeHandler(service))
		necessidade.GET("", handlers.ListNecessidadesHandler(service))
		necessidade.GET(":id", handlers.GetNecessidadeHandler(service))
		necessidade.PATCH(":id", handlers.UpdateNecessidadeHandler(service))
		necessidade.GET(":id/matches", handlers.GetMatchesHandler(matchingService))
		necessidade.GET(":id/desafios", handlers.ListDesafiosByNecessidadeHandler(desafioService))
	}
}
