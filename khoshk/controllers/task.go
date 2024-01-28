package controllers

import (
	"PlanVerse/configs"
	"PlanVerse/messages"
	"PlanVerse/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func CreateTaskHandler(ctx echo.Context) error {
	req := new(models.CreateTaskRequest)
	res := new(models.CreateTaskResponse)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	var stateIDs []int
	result := configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	exist := false
	for _, id := range stateIDs {
		if id == req.StateID {
			exist = true
			break
		}
	}
	if !exist {
		return ctx.JSON(http.StatusNotAcceptable, messages.NotInProject)
	}
	var state models.State
	result = configs.DB.Where("id = ?", req.StateID).Preload("Tasks").Find(&state)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongStateID)
	}
	newTask := models.Task{
		Title:           req.Title,
		BackGroundColor: req.BackGroundColor,
		Description:     req.Description,
	}
	state.Tasks = append(state.Tasks, newTask)
	result = configs.DB.Save(&state)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	res.TaskID = int(state.Tasks[len(state.Tasks)-1].ID)
	return ctx.JSON(http.StatusOK, res)
}
