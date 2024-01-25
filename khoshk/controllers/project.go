package controllers

import (
	"PlanVerse/configs"
	"PlanVerse/messages"
	"PlanVerse/models"
	_ "PlanVerse/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ProjectListHandler(ctx echo.Context) error {
	var projectsList []models.ResponseList
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	result := configs.DB.Table("projects").Select([]string{"projects.id", "projects.title", "projects.back_ground_pic", "projects.members_number"}).Joins("join projects_members on projects_members.project_id = projects.id").Where("user_id = ?", userID).Scan(&projectsList)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	return ctx.JSON(http.StatusOK, projectsList)
}
