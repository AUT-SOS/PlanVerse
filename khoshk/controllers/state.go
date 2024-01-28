package controllers

import (
	"PlanVerse/configs"
	"PlanVerse/helpers"
	"PlanVerse/messages"
	"PlanVerse/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func StateListHandler(ctx echo.Context) error {
	var res []models.StateObject
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	var projectIDs []helpers.ProjectID
	result := configs.DB.Table("projects_members").Select("project_id").Where("user_id = ?", userID).Scan(&projectIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	exist := false
	for _, id := range projectIDs {
		if id.ProjectID == projectID {
			exist = true
			break
		}
	}
	if !exist {
		return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
	}
	result = configs.DB.Table("states").Select([]string{"title", "back_ground_color", "admin_access"}).Where("project_id = ?", projectID).Scan(&res)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	return ctx.JSON(http.StatusOK, res)
}

func CreateStateHandler(ctx echo.Context) error {
	req := new(models.CreateStateRequest)
	res := new(models.CreateStateResponse)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	var project models.Project
	result := configs.DB.Where("id = ?", projectID).Preload("States").Find(&project)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	newState := models.State{
		Title:           req.Title,
		BackGroundColor: req.BackGroundColor,
		AdminAccess:     req.AdminAccess,
	}
	project.States = append(project.States, newState)
	result = configs.DB.Save(&project)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	res.StateID = int(project.States[len(project.States)-1].ID)
	return ctx.JSON(http.StatusOK, res)
}

func EditStateHandler(ctx echo.Context) error {
	req := new(models.StateObject)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	var project models.Project
	result := configs.DB.Where("id = ?", projectID).Preload("States").Find(&project)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	for i := 0; i < len(project.States); i++ {
		if int(project.States[i].ID) == req.ID {
			project.States[i].Title = req.Title
			project.States[i].BackGroundColor = req.BackGroundColor
			project.States[i].AdminAccess = req.AdminAccess
			break
		}
	}
	result = configs.DB.Save(&project)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	return ctx.JSON(http.StatusOK, messages.StateEdited)
}
