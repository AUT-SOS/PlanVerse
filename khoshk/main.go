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
	apiGroup := server.Group("")

	//middlewares
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
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
	userApiGroup := apiGroup.Group("/user")
	userApiGroup.POST("/register", controllers.RegisterHandler)
	userApiGroup.POST("/verify", controllers.VerifyHandler)
	userApiGroup.POST("/refresh", controllers.RefreshHandler)
	userApiGroup.POST("/login", controllers.LoginHandler)
	userApiGroup.POST("/resend-email", controllers.ResendEmailHandler)
	userApiGroup.GET("/get-user/:user-id", controllers.GetUserHandler)
	userApiGroup.GET("/get-my-user", controllers.GetUserIDHandler)

	//project api
	projectApiGroup := apiGroup.Group("/project")
	projectApiGroup.GET("/list", controllers.ProjectListHandler)

	log.Fatal(server.Start("localhost:8080"))
}
