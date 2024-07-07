package main

import (
	"PlanVerse/configs"
	"PlanVerse/models"
	"PlanVerse/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
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
	customRegistry := prometheus.NewRegistry()
	customRegistry.MustRegister(models.SuccessRequests)
	customRegistry.MustRegister(models.FailedRequests)
	customRegistry.MustRegister(models.SuccessDBRequests)
	customRegistry.MustRegister(models.FailedDBRequests)
	customRegistry.MustRegister(models.ResponseTime)
	configs.ConnectToDatabase()
	configs.ConnectToRedis()
	server := echo.New()
	routes.Routes(server, customRegistry)
}
