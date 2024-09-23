package main

import (
	"golang-ai-management/models"
	routes "golang-ai-management/routes"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func main() {
	var serverConfig = models.Config

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	//router.Use(middleware.Authentication())

	routes.FaceRegRoutes(router)

	slog.Info("Server is running in port %d", serverConfig.Port)
	router.Run(":" + serverConfig.Port)
}
