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
		return ctx.JSON(http.StatusNotAcceptable, messages.TaskNotInProject)
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

func ChangeTaskStateHandler(ctx echo.Context) error {
	req := new(models.ChangeTaskStateRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	var sourceStateID int
	result := configs.DB.Table("tasks").Select("state_id").Where("id = ?", req.TaskID).Scan(&sourceStateID)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongTaskID)
	}
	var stateIDs []int
	result = configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	srcExist := false
	desExist := false
	for _, id := range stateIDs {
		if id == req.StateID {
			desExist = true
		} else if id == sourceStateID {
			srcExist = true
		}
	}
	if !srcExist || !desExist {
		return ctx.JSON(http.StatusNotAcceptable, messages.TaskNotInProject)
	}
	var srcAdminAccess bool
	result = configs.DB.Table("states").Select("admin_access").Where("id = ?", sourceStateID).Scan(&srcAdminAccess)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var desAdminAccess bool
	result = configs.DB.Table("states").Select("admin_access").Where("id = ?", req.StateID).Scan(&desAdminAccess)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if srcAdminAccess || desAdminAccess {
		userIDCtx := ctx.Get("user_id")
		userID := userIDCtx.(int)
		var isAdmin bool
		result = configs.DB.Table("projects_members").Select("is_admin").Where("project_id = ? and user_id = ?", projectID, userID).Scan(&isAdmin)
		if result.Error != nil {
			return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
		}
		if !isAdmin {
			return ctx.JSON(http.StatusUnauthorized, messages.AdminChange)
		}
		result = configs.DB.Table("tasks").Where("id = ?", req.TaskID).Update("state_id", req.StateID)
		if result.Error != nil {
			return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
		}
	}
	result = configs.DB.Table("tasks").Where("id = ?", req.TaskID).Update("state_id", req.StateID)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	return ctx.JSON(http.StatusOK, messages.TaskStateChanged)
}

func AddPerformerHandler(ctx echo.Context) error {
	req := new(models.PerformerRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	var stateID int
	result := configs.DB.Table("tasks").Select("state_id").Where("id = ?", req.TaskID).Scan(&stateID)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongTaskID)
	}
	var stateIDs []int
	result = configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	exist := false
	for _, id := range stateIDs {
		if id == stateID {
			exist = true
			break
		}
	}
	if !exist {
		return ctx.JSON(http.StatusNotAcceptable, messages.TaskNotInProject)
	}
	var members []int
	result = configs.DB.Table("projects_members").Select("user_id").Where("project_id = ?", projectID).Scan(&members)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	exist = false
	for _, id := range members {
		if id == req.PerformerID {
			exist = true
			break
		}
	}
	if !exist {
		return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
	}
	var performerIDs []int
	result = configs.DB.Table("tasks_performers").Select("user_id").Where("task_id = ?", req.TaskID).Scan(&performerIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongTaskID)
	}
	exist = false
	for _, id := range performerIDs {
		if id == req.PerformerID {
			exist = true
			break
		}
	}
	if exist {
		return ctx.JSON(http.StatusNotAcceptable, messages.InTask)
	}
	var user models.User
	result = configs.DB.Where("id = ?", req.PerformerID).Find(&user)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongUserID)
	}
	var task models.Task
	result = configs.DB.Where("id = ?", req.TaskID).Preload("Performers").Find(&task)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	task.Performers = append(task.Performers, user)
	result = configs.DB.Save(&task)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	return ctx.JSON(http.StatusOK, messages.TaskAssigned)
}

func RemovePerformerHandler(ctx echo.Context) error {
	req := new(models.PerformerRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	var stateID int
	result := configs.DB.Table("tasks").Select("state_id").Where("id = ?", req.TaskID).Scan(&stateID)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongTaskID)
	}
	var stateIDs []int
	result = configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	exist := false
	for _, id := range stateIDs {
		if id == stateID {
			exist = true
			break
		}
	}
	if !exist {
		return ctx.JSON(http.StatusNotAcceptable, messages.TaskNotInProject)
	}
	var members []int
	result = configs.DB.Table("projects_members").Select("user_id").Where("project_id = ?", projectID).Scan(&members)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	exist = false
	for _, id := range members {
		if id == req.PerformerID {
			exist = true
			break
		}
	}
	if !exist {
		return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
	}
	var performerIDs []int
	result = configs.DB.Table("tasks_performers").Select("user_id").Where("task_id = ?", req.TaskID).Scan(&performerIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongTaskID)
	}
	exist = false
	for _, id := range performerIDs {
		if id == req.PerformerID {
			exist = true
			break
		}
	}
	if !exist {
		return ctx.JSON(http.StatusNotAcceptable, messages.NotInTask)
	}
	result = configs.DB.Unscoped().Where("task_id = ? and user_id = ?", req.TaskID, req.PerformerID).Delete(&models.TasksPerformers{})
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	return ctx.JSON(http.StatusOK, messages.PerformerRemoved)
}

func EditTaskHandler(ctx echo.Context) error {
	req := new(models.EditTaskRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	var stateID int
	result := configs.DB.Table("tasks").Select("state_id").Where("id = ?", req.TaskID).Scan(&stateID)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongTaskID)
	}
	var stateIDs []int
	result = configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	exist := false
	for _, id := range stateIDs {
		if id == stateID {
			exist = true
			break
		}
	}
	if !exist {
		return ctx.JSON(http.StatusNotAcceptable, messages.TaskNotInProject)
	}
	var task models.Task
	result = configs.DB.Where("id = ?", req.TaskID).Find(&task)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	task.Title = req.Title
	task.BackGroundColor = req.BackGroundColor
	task.Description = req.Description
	result = configs.DB.Save(&task)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	return ctx.JSON(http.StatusOK, messages.TaskEdited)
}

func DeleteTaskHandler(ctx echo.Context) error {
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	taskID, err := strconv.Atoi(ctx.Param("task-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongTaskID)
	}
	var stateID int
	result := configs.DB.Table("tasks").Select("state_id").Where("id = ?", taskID).Scan(&stateID)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongTaskID)
	}
	var stateIDs []int
	result = configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	exist := false
	for _, id := range stateIDs {
		if id == stateID {
			exist = true
			break
		}
	}
	if !exist {
		return ctx.JSON(http.StatusNotAcceptable, messages.TaskNotInProject)
	}
	result = configs.DB.Unscoped().Where("task_id = ?", taskID).Delete(&models.TasksPerformers{})
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	result = configs.DB.Unscoped().Where("id = ?", taskID).Delete(&models.Task{})
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	return ctx.JSON(http.StatusOK, messages.TaskDeleted)
}
