package routes

import (
	"conectainova/internal/handlers"
	"conectainova/internal/services"
	"github.com/gin-gonic/gin"
)

func initLojistaRoutes(rg *gin.RouterGroup, service services.ILojistaService) {
	lojista := rg.Group("lojistas")
	{
		lojista.POST("", handlers.CreateLojistaHandler(service))
		lojista.GET(":id", handlers.GetLojistaHandler(service))
		lojista.PATCH(":id", handlers.UpdateLojistaHandler(service))
	}
}
