package routes

import (
	"conectainova/internal/handlers"
	"conectainova/internal/services"
	"github.com/gin-gonic/gin"
)

func initDesafioRoutes(rg *gin.RouterGroup, service services.IDesafioService) {
	desafio := rg.Group("desafios")
	{
		desafio.POST("", handlers.CreateDesafioHandler(service))
		desafio.GET("", handlers.ListDesafiosHandler(service))
		desafio.PATCH(":id", handlers.UpdateDesafioHandler(service))
	}
}
