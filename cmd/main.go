package main

import (
	"fmt"
	main2 "jumyste-app-backend/config"
	"jumyste-app-backend/internal/applicator"
	"jumyste-app-backend/internal/middleware"
	"jumyste-app-backend/internal/router"
	"jumyste-app-backend/pkg/logger"
)

func main() {
	logger.InitLogger()

	main2.LoadConfig()

	logger.Log.Info("Starting application...")
	app := applicator.NewApp()

	auth := middleware.NewAuthMiddleware(main2.AppConfig)

	r := router.SetupRouter(app.AuthHandler, app.UserHandler, auth)

	serverPort := main2.AppConfig.Server.Port
	addr := fmt.Sprintf(":%s", serverPort)
	logger.Log.Info("Starting server", "port", serverPort)

	if err := r.Run(addr); err != nil {
		logger.Log.Error("Failed to start server", "error", err.Error())
	}
}
