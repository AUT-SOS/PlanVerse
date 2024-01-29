package controllers

import (
	"PlanVerse/configs"
	"PlanVerse/helpers"
	"PlanVerse/messages"
	"PlanVerse/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"sync"
)

func StateListHandler(ctx echo.Context) error {
	var res []models.StateListResponse
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
	result = configs.DB.Table("states").Select([]string{"id", "title", "back_ground_color", "admin_access"}).Where("project_id = ?", projectID).Scan(&res)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	var wg sync.WaitGroup
	for i := range res {
		wg.Add(1)
		go func(index int, wg *sync.WaitGroup) {
			defer wg.Done()
			var taskShows []models.TaskShow
			result = configs.DB.Table("tasks").Select([]string{"id", "title", "back_ground_color"}).Where("state_id = ?", res[index].ID).Scan(&taskShows)
			for j, task := range taskShows {
				var performers []int
				result = configs.DB.Table("tasks_performers").Select("user_id").Where("task_id = ?", task.ID).Scan(&performers)
				taskShows[j].Performers = performers
			}
			res[index].Tasks = taskShows
		}(i, &wg)
	}
	wg.Wait()
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
	var stateIDs []int
	result := configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	exist := false
	for _, id := range stateIDs {
		if id == req.ID {
			exist = true
			break
		}
	}
	if !exist {
		return ctx.JSON(http.StatusNotAcceptable, messages.StateNotInProject)
	}
	var state models.State
	result = configs.DB.Where("id = ?", req.ID).Find(&state)
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
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	stateID, err := strconv.Atoi(ctx.Param("state-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongStateID)
	}
	var stateIDs []int
	result := configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
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
		return ctx.JSON(http.StatusNotAcceptable, messages.StateNotInProject)
	}
	var taskIDs []int
	result = configs.DB.Table("tasks").Select("id").Where("state_id = ?", stateID).Scan(&taskIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var wg *sync.WaitGroup
	for i := range taskIDs {
		wg.Add(1)
		go func(index int, wg *sync.WaitGroup) {
			defer wg.Done()
			configs.DB.Unscoped().Where("task_id = ?", taskIDs[index]).Delete(&models.TasksPerformers{})
		}(i, wg)
	}
	wg.Wait()
	result = configs.DB.Unscoped().Where("state_id = ?", stateID).Delete(&models.Task{})
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	result = configs.DB.Unscoped().Where("id = ?", stateID).Delete(&models.State{})
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
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
	var taskIDs []int
	result = configs.DB.Table("tasks").Select("id").Where("state_id = ?", stateID).Scan(&taskIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := range taskIDs {
		wg.Add(1)
		go func(index int, wg *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()
			var task models.GetTaskResponse
			configs.DB.Table("tasks").Select([]string{"title", "back_ground_color", "description"}).Where("id = ?", taskIDs[index]).Scan(&task)
			var performerIDs []int
			configs.DB.Table("tasks_performers").Select("user_id").Where("task_id = ?", taskIDs[index]).Scan(&performerIDs)
			task.ID = taskIDs[index]
			task.Performers = performerIDs
			mu.Lock()
			res.Tasks = append(res.Tasks, task)
			mu.Unlock()
		}(i, &wg, &mu)
	}
	res.ID = stateID
	res.ProjectID = projectID.ProjectID
	wg.Wait()
	return ctx.JSON(http.StatusOK, res)
}
