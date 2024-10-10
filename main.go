package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-ai-management/config"
	"golang-ai-management/logger"
	routes "golang-ai-management/routes"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
	var serverConfig = config.Config

	factory := logger.LoggerFactory{}
	newLogger, err := factory.NewLogger(serverConfig.LogType, serverConfig.LogFormat, serverConfig.LogLevel)

	if err != nil {
		fmt.Println("Error creating newLogger:", err)
		return
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	//router.Use(middleware.Authentication())

	routes.FaceRegRoutes(router)
	newLogger.InfoArgs("Server running in port " + serverConfig.Port)

	err = router.Run(":" + serverConfig.Port)
	if err != nil {
		return
	}

}
