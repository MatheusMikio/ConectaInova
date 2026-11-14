package routes

import (
	"net/http"

	"conectainova/config"
	"conectainova/internal/container"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(db *gorm.DB) error {
	container := container.NewContainer(db)
	router := gin.Default()
	router.Use(corsWithCredentials(config.GetCORSAllowedOrigins()))
	initRoutes(router, container)

	return router.Run(":8080")
}

func corsWithCredentials(allowedOrigins []string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowed[origin] = struct{}{}
	}

	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		if _, ok := allowed[origin]; ok {
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Access-Control-Allow-Credentials", "true")
			ctx.Header("Vary", "Origin")
		}

		if ctx.Request.Method == http.MethodOptions {
			ctx.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
			ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			ctx.Status(http.StatusNoContent)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
