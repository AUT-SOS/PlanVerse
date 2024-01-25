package main

import (
	"PlanVerse/configs"
	"PlanVerse/controllers"
	"PlanVerse/middlewares"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if os.Getenv("ENV") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("App .env file not found")
		}
	}
	configs.ConnectToDatabase()
	configs.ConnectToRedis()
	server := echo.New()
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))
	server.Use(middlewares.AuthMiddleware)
	server.POST("/register", controllers.RegisterHandler)
	server.POST("/verify", controllers.VerifyHandler)
	server.POST("/refresh", controllers.RefreshHandler)
	server.POST("/login", controllers.LoginHandler)
	server.POST("/resend-email", controllers.ResendEmailHandler)
	server.GET("/get-user/:user-id", controllers.GetUserHandler)
	server.GET("/get-my-user", controllers.GetUserIDHandler)
	log.Fatal(server.Start("localhost:8080"))
}
