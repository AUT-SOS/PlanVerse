package main

import (
	"PlanVerse/configs"
	"PlanVerse/controllers"
	"PlanVerse/middlewares"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"os"
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
	server.Use(middlewares.AuthMiddleware)
	server.POST("/register", controllers.RegisterHandler)
	server.POST("/verify", controllers.VerifyHandler)
	server.POST("/refresh", controllers.RefreshHandler)
	server.POST("/login", controllers.LoginHandler)
	server.GET("/get-user", controllers.GetUserHandler)
	log.Fatal(server.Start("localhost:8080"))
}
