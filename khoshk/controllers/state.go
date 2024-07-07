package controllers

import (
	"PlanVerse/configs"
	"PlanVerse/messages"
	"PlanVerse/models"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

func StateListHandler(ctx echo.Context) error {
	startTime := time.Now()
	method := ctx.Request().Method
	endpoint := ctx.Request().URL.Path
	var res []models.StateListResponse
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	var projectIDs []int
	result := configs.DB.Table("projects_members").Select("project_id").Where("user_id = ?", userID).Scan(&projectIDs)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	if len(projectIDs) == 0 {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongUserID)
	}
	exist := false
	for _, id := range projectIDs {
		if id == projectID {
			exist = true
			break
		}
	}
	if !exist {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
	}
	result = configs.DB.Table("states").Select([]string{"id", "title", "back_ground_color", "admin_access"}).Where("project_id = ?", projectID).Scan(&res)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	var wg sync.WaitGroup
	for i := range res {
		wg.Add(1)
		go func(index int, wg *sync.WaitGroup) {
			defer wg.Done()
			var taskShows []models.TaskShow
			configs.DB.Table("tasks").Select([]string{"id", "index", "title", "back_ground_color", "description", "deadline", "estimated_time", "actual_time", "priority"}).Where("state_id = ?", res[index].ID).Scan(&taskShows)
			for j, task := range taskShows {
				var performers []int
				configs.DB.Table("tasks_performers").Select("user_id").Where("task_id = ?", task.ID).Scan(&performers)
				taskShows[j].Performers = performers
			}
			res[index].Tasks = taskShows
		}(i, &wg)
	}
	wg.Wait()
	models.SuccessRequests.WithLabelValues(method, endpoint).Inc()
	models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
	return ctx.JSON(http.StatusOK, res)
}

func CreateStateHandler(ctx echo.Context) error {
	startTime := time.Now()
	method := ctx.Request().Method
	endpoint := ctx.Request().URL.Path
	req := new(models.CreateStateRequest)
	res := new(models.CreateStateResponse)
	if err := ctx.Bind(req); err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	var project models.Project
	result := configs.DB.Where("id = ?", projectID).Preload("States").Find(&project)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	if project.ID == 0 {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
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
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	res.StateID = int(project.States[len(project.States)-1].ID)
	models.SuccessRequests.WithLabelValues(method, endpoint).Inc()
	models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
	return ctx.JSON(http.StatusOK, res)
}

func EditStateHandler(ctx echo.Context) error {
	startTime := time.Now()
	method := ctx.Request().Method
	endpoint := ctx.Request().URL.Path
	req := new(models.StateObject)
	if err := ctx.Bind(req); err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	stateID, err := strconv.Atoi(ctx.Param("state-id"))
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongStateID)
	}
	var stateIDs []int
	result := configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	if len(stateIDs) == 0 {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
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
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusNotAcceptable, messages.StateNotInProject)
	}
	var state models.State
	result = configs.DB.Where("id = ?", stateID).Find(&state)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	state.Title = req.Title
	state.BackGroundColor = req.BackGroundColor
	state.AdminAccess = req.AdminAccess
	result = configs.DB.Save(&state)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	models.SuccessRequests.WithLabelValues(method, endpoint).Inc()
	models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
	return ctx.JSON(http.StatusOK, messages.StateEdited)
}

func DeleteStateHandler(ctx echo.Context) error {
	startTime := time.Now()
	method := ctx.Request().Method
	endpoint := ctx.Request().URL.Path
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	stateID, err := strconv.Atoi(ctx.Param("state-id"))
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongStateID)
	}
	var stateIDs []int
	result := configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if len(stateIDs) == 0 {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
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
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusNotAcceptable, messages.StateNotInProject)
	}
	var taskIDs []int
	result = configs.DB.Table("tasks").Select("id").Where("state_id = ?", stateID).Scan(&taskIDs)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	var wg sync.WaitGroup
	for i := range taskIDs {
		wg.Add(1)
		go func(index int, wg *sync.WaitGroup) {
			defer wg.Done()
			configs.DB.Unscoped().Where("task_id = ?", taskIDs[index]).Delete(&models.TasksPerformers{})
		}(i, &wg)
	}
	wg.Wait()
	result = configs.DB.Unscoped().Where("state_id = ?", stateID).Delete(&models.Task{})
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	result = configs.DB.Unscoped().Where("id = ?", stateID).Delete(&models.State{})
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	models.SuccessRequests.WithLabelValues(method, endpoint).Inc()
	models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
	return ctx.JSON(http.StatusOK, messages.StateDeleted)
}

func GetStateHandler(ctx echo.Context) error {
	startTime := time.Now()
	method := ctx.Request().Method
	endpoint := ctx.Request().URL.Path
	res := new(models.GetStateResponse)
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	stateID, err := strconv.Atoi(ctx.Param("state-id"))
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongStateID)
	}
	var stateIDs []int
	result := configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	if len(stateIDs) == 0 {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
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
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusNotAcceptable, messages.StateNotInProject)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	var projectIDs []int
	result = configs.DB.Table("projects_members").Select("project_id").Where("user_id = ?", userID).Scan(&projectIDs)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	if len(projectIDs) == 0 {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusNotAcceptable, messages.UserNoProject)
	}
	exist = false
	for _, id := range projectIDs {
		if id == projectID {
			exist = true
			break
		}
	}
	if !exist {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
	}
	result = configs.DB.Table("states").Select([]string{"title", "back_ground_color", "admin_access"}).Where("id = ?", stateID).Scan(res)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	var taskIDs []int
	result = configs.DB.Table("tasks").Select("id").Where("state_id = ?", stateID).Scan(&taskIDs)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := range taskIDs {
		wg.Add(1)
		go func(index int, wg *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()
			var task models.GetTaskResponse
			configs.DB.Table("tasks").Select([]string{"index", "title", "back_ground_color", "description", "deadline", "estimated_time", "actual_time", "priority"}).Where("id = ?", taskIDs[index]).Scan(&task)
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
	wg.Wait()
	models.SuccessRequests.WithLabelValues(method, endpoint).Inc()
	models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
	return ctx.JSON(http.StatusOK, res)
}
