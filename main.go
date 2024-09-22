package main

import (
	"log/slog"
	"os"

	routes "golang-ai-management/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	//router.Use(middleware.Authentication())

	routes.FaceRegRoutes(router)
	//routes.MenuRoutes(router)
	//routes.TableRoutes(router)
	//routes.OrderRoutes(router)
	//routes.OrderItemRoutes(router)
	//routes.InvoiceRoutes(router)

	slog.Info("Server is running in port 8000")
	router.Run(":" + port)
}
