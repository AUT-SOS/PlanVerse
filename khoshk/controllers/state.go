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
	var state models.State
	result := configs.DB.Where("id = ?", req.ID).Find(&state)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongStateID)
	}
	state.Title = req.Title
	state.BackGroundColor = req.BackGroundColor
	state.AdminAccess = req.AdminAccess
	result = configs.DB.Save(&state)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	return ctx.JSON(http.StatusOK, messages.StateEdited)
}

func DeleteStateHandler(ctx echo.Context) error {
	stateID, err := strconv.Atoi(ctx.Param("state-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongStateID)
	}
	result := configs.DB.Unscoped().Where("id = ?", stateID).Delete(&models.State{})
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongStateID)
	}
	return ctx.JSON(http.StatusOK, messages.StateDeleted)
}

func GetStateHandler(ctx echo.Context) error {
	res := new(models.GetStateResponse)
	stateID, err := strconv.Atoi(ctx.Param("state-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongStateID)
	}
	var projectID helpers.ProjectID
	result := configs.DB.Table("states").Select("project_id").Where("id = ?", stateID).Scan(&projectID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongStateID)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	var projectIDs []helpers.ProjectID
	result = configs.DB.Table("projects_members").Select("project_id").Where("user_id = ?", userID).Scan(&projectIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	exist := false
	for _, id := range projectIDs {
		if id.ProjectID == projectID.ProjectID {
			exist = true
			break
		}
	}
	if !exist {
		return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
	}
	result = configs.DB.Table("states").Select([]string{"title", "back_ground_color", "admin_access"}).Where("id = ?", stateID).Scan(res)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	res.ID = stateID
	res.ProjectID = projectID.ProjectID
	return ctx.JSON(http.StatusOK, res)
}
