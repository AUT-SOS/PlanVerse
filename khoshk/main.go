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

	//middlewares
	//server.Use(middleware.Logger())
	//server.Use(middleware.Recover())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))
	server.Use(middlewares.AuthMiddleware)
	server.Use(middlewares.VerifyMiddleware)

	//user api
	server.POST("/register", controllers.RegisterHandler)
	server.POST("/verify", controllers.VerifyHandler)
	server.POST("/refresh", controllers.RefreshHandler)
	server.POST("/login", controllers.LoginHandler)
	server.POST("/resend-email", controllers.ResendEmailHandler)
	server.GET("/get-user/:user-id", controllers.GetUserHandler)
	server.GET("/get-my-user", controllers.GetUserIDHandler)

	//project api
	server.GET("/list", controllers.ProjectListHandler)
	server.POST("/create", controllers.CreateProjectHandler)
	server.POST("/share-link", controllers.ShareProjectHandler)

	//start server
	log.Fatal(server.Start("localhost:8080"))
}
