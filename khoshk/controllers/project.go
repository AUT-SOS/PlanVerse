package controllers

import (
	"PlanVerse/configs"
	"PlanVerse/helpers"
	"PlanVerse/messages"
	"PlanVerse/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ProjectListHandler(ctx echo.Context) error {
	var projectsList []models.ProjectListResponse
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	result := configs.DB.Table("projects").Select([]string{"projects.id", "projects.title", "projects.back_ground_pic", "projects.members_number", "projects_members.is_admin"}).Joins("inner join projects_members on projects.id = projects_members.project_id").Where("projects_members.user_id = ?", userID).Scan(&projectsList)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	return ctx.JSON(http.StatusOK, projectsList)
}

func CreateProjectHandler(ctx echo.Context) error {
	req := new(models.CreateProjectRequest)
	res := new(models.CreateProjectResponse)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	var user models.User
	result := configs.DB.Where("id = ?", userID).Find(&user)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	link, err := helpers.GenerateRandomLink(req.Title)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToGenerateLink)
	}
	joinLink := models.JoinLink{
		Link: link,
	}
	newProject := models.Project{
		Title:         req.Title,
		Description:   req.Description,
		BackGroundPic: req.Picture,
		OwnerID:       userID,
		MembersNumber: 1,
		Members:       []models.User{user},
		JoinLink:      joinLink,
	}
	if err = configs.DB.Create(&newProject).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToCreateProject)
	}
	result = configs.DB.Table("projects_members").Where("project_id = ?", newProject.ID).Update("is_admin", true)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	res.ProjectID = int(newProject.ID)
	res.Link = joinLink.Link
	res.Message = messages.ProjectCreated
	return ctx.JSON(http.StatusOK, res)
}
