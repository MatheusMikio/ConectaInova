package routes

import (
	"conectainova/internal/handlers"
	"conectainova/internal/services"
	"github.com/gin-gonic/gin"
)

func initStartupRoutes(rg *gin.RouterGroup, service services.IStartupService) {
	startup := rg.Group("startups")
	{
		startup.POST("", handlers.CreateStartupHandler(service))
		startup.GET("", handlers.ListStartupsHandler(service))
		startup.PATCH(":id", handlers.UpdateStartupHandler(service))
	}
}
