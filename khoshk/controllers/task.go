package controllers

import (
	"PlanVerse/configs"
	"PlanVerse/messages"
	"PlanVerse/models"
	"PlanVerse/services"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CreateTaskHandler(ctx echo.Context) error {
	startTime := time.Now()
	method := ctx.Request().Method
	endpoint := ctx.Request().URL.Path
	req := new(models.CreateTaskRequest)
	res := new(models.CreateTaskResponse)
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
	result = configs.DB.Where("id = ?", stateID).Preload("Tasks").Find(&state)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	if state.ID == 0 {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongStateID)
	}
	if req.EstimatedTime == 0 {
		req.EstimatedTime = 10
	}
	if req.Deadline == "" {
		endTime := time.Now().Add(10 * 24 * time.Hour)
		deadline := endTime.String()
		req.Deadline = deadline[:10]
	}
	if req.Priority < 1 || req.Priority > 5 {
		req.Priority = 5
	}
	year, _ := strconv.Atoi(req.Deadline[:4])
	month, _ := strconv.Atoi(req.Deadline[5:7])
	day, _ := strconv.Atoi(req.Deadline[8:])
	newTask := models.Task{
		Title:           req.Title,
		BackGroundColor: req.BackGroundColor,
		Description:     req.Description,
		EstimatedTime:   req.EstimatedTime,
		ActualTime:      0,
		Deadline:        time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC),
		Priority:        req.Priority,
		Index:           req.Index,
	}
	state.Tasks = append(state.Tasks, newTask)
	result = configs.DB.Save(&state)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	res.TaskID = int(state.Tasks[len(state.Tasks)-1].ID)
	go func() {
		time.Sleep(24 * time.Hour)
		var task models.Task
		_ = configs.DB.Where("id = ?", res.TaskID).Find(&task)
		var taskState models.State
		_ = configs.DB.Where("id = ?", task.StateID).Find(&taskState)
		if taskState.Title == "Done" {
			return
		}
		task.ActualTime = task.ActualTime + 1
		_ = configs.DB.Save(&task)
	}()
	var members []int
	result = configs.DB.Table("projects_members").Select("user_id").Where("project_id = ?", projectID).Scan(&members)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	var task models.Task
	result = configs.DB.Where("id = ?", res.TaskID).Preload("Performers").Find(&task)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	wsMessage := models.WSMessage{
		Type:    "update-task",
		Payload: task,
	}
	for conn, client := range models.Clients {
		for i := range members {
			if client == members[i] {
				wsErr := conn.WriteJSON(wsMessage)
				if wsErr != nil {
					delete(models.Clients, conn)
				}
			}
		}
	}
	models.SuccessRequests.WithLabelValues(method, endpoint).Inc()
	models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
	return ctx.JSON(http.StatusOK, res)
}

func ChangeTaskStateHandler(ctx echo.Context) error {
	startTime := time.Now()
	method := ctx.Request().Method
	endpoint := ctx.Request().URL.Path
	req := new(models.ChangeTaskStateRequest)
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
	taskID, err := strconv.Atoi(ctx.Param("task-id"))
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongTaskID)
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
	var sourceStateID int
	result = configs.DB.Table("tasks").Select("state_id").Where("id = ?", taskID).Scan(&sourceStateID)
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	if sourceStateID == 0 {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongTaskID)
	}
	var stateIDs []int
	result = configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
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
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusNotAcceptable, messages.TaskNotInProject)
	}
	var srcAdminAccess bool
	result = configs.DB.Table("states").Select("admin_access").Where("id = ?", sourceStateID).Scan(&srcAdminAccess)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	var desAdminAccess bool
	result = configs.DB.Table("states").Select("admin_access").Where("id = ?", req.StateID).Scan(&desAdminAccess)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	if srcAdminAccess || desAdminAccess {
		var isAdmin bool
		result = configs.DB.Table("projects_members").Select("is_admin").Where("project_id = ? and user_id = ?", projectID, userID).Scan(&isAdmin)
		if result.Error != nil {
			models.FailedRequests.WithLabelValues(method, endpoint).Inc()
			models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
			models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
			return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
		}
		models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
		if !isAdmin {
			models.FailedRequests.WithLabelValues(method, endpoint).Inc()
			models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
			return ctx.JSON(http.StatusUnauthorized, messages.AdminChange)
		}
		result = configs.DB.Table("tasks").Where("id = ?", taskID).Update("state_id", req.StateID)
		if result.Error != nil {
			models.FailedRequests.WithLabelValues(method, endpoint).Inc()
			models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
			models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
			return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
		}
		models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
		models.SuccessRequests.WithLabelValues(method, endpoint).Inc()
		return ctx.JSON(http.StatusOK, messages.TaskStateChanged)
	}
	result = configs.DB.Table("tasks").Where("id = ?", taskID).Update("state_id", req.StateID)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	var members []int
	result = configs.DB.Table("projects_members").Select("user_id").Where("project_id = ?", projectID).Scan(&members)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	var task models.Task
	result = configs.DB.Where("id = ?", taskID).Preload("Performers").Find(&task)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	wsMessage := models.WSMessage{
		Type:    "update-task",
		Payload: task,
	}
	for conn, client := range models.Clients {
		for i := range members {
			if client == members[i] {
				wsErr := conn.WriteJSON(wsMessage)
				if wsErr != nil {
					delete(models.Clients, conn)
				}
			}
		}
	}
	models.SuccessRequests.WithLabelValues(method, endpoint).Inc()
	models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
	return ctx.JSON(http.StatusOK, messages.TaskStateChanged)
}

func AddPerformerHandler(ctx echo.Context) error {
	startTime := time.Now()
	method := ctx.Request().Method
	endpoint := ctx.Request().URL.Path
	req := new(models.PerformerRequest)
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
	taskID, err := strconv.Atoi(ctx.Param("task-id"))
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongTaskID)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	var stateID int
	result := configs.DB.Table("tasks").Select("state_id").Where("id = ?", taskID).Scan(&stateID)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	if stateID == 0 {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongTaskID)
	}
	var stateIDs []int
	result = configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
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
		return ctx.JSON(http.StatusNotAcceptable, messages.TaskNotInProject)
	}
	var members []int
	result = configs.DB.Table("projects_members").Select("user_id").Where("project_id = ?", projectID).Scan(&members)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	exist = false
	for _, id := range members {
		if id == req.PerformerID {
			exist = true
			break
		}
	}
	if !exist {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
	}
	var performerIDs []int
	result = configs.DB.Table("tasks_performers").Select("user_id").Where("task_id = ?", taskID).Scan(&performerIDs)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	exist = false
	for _, id := range performerIDs {
		if id == req.PerformerID {
			exist = true
			break
		}
	}
	if exist {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusNotAcceptable, messages.InTask)
	}
	var user models.User
	result = configs.DB.Where("id = ?", req.PerformerID).Find(&user)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	if user.ID == 0 {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongUserID)
	}
	var task models.Task
	result = configs.DB.Where("id = ?", taskID).Preload("Performers").Find(&task)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	task.Performers = append(task.Performers, user)
	result = configs.DB.Save(&task)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	var email string
	result = configs.DB.Table("users").Select("email").Where("id = ?", req.PerformerID).Scan(&email)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var projectTitle string
	result = configs.DB.Table("projects").Select("title").Where("id = ?", projectID).Scan(&projectTitle)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	var username string
	result = configs.DB.Table("users").Select("username").Where("id = ?", userID).Scan(&username)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	go func() {
		services.SendMail("PlanVerse Notification", fmt.Sprintf("you've been assigned to new task in %s project by %s!", projectTitle, username), []string{email})
	}()
	wsMessage := models.WSMessage{
		Type:    "update-task",
		Payload: task,
	}
	for conn, client := range models.Clients {
		for i := range members {
			if client == members[i] {
				wsErr := conn.WriteJSON(wsMessage)
				if wsErr != nil {
					delete(models.Clients, conn)
				}
			}
		}
	}
	models.SuccessRequests.WithLabelValues(method, endpoint).Inc()
	models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
	return ctx.JSON(http.StatusOK, messages.TaskAssigned)
}

func RemovePerformerHandler(ctx echo.Context) error {
	startTime := time.Now()
	method := ctx.Request().Method
	endpoint := ctx.Request().URL.Path
	req := new(models.PerformerRequest)
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
	taskID, err := strconv.Atoi(ctx.Param("task-id"))
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongTaskID)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	var stateID int
	result := configs.DB.Table("tasks").Select("state_id").Where("id = ?", taskID).Scan(&stateID)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	if stateID == 0 {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongTaskID)
	}
	var stateIDs []int
	result = configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
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
		return ctx.JSON(http.StatusNotAcceptable, messages.TaskNotInProject)
	}
	var members []int
	result = configs.DB.Table("projects_members").Select("user_id").Where("project_id = ?", projectID).Scan(&members)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	exist = false
	for _, id := range members {
		if id == req.PerformerID {
			exist = true
			break
		}
	}
	if !exist {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
	}
	var performerIDs []int
	result = configs.DB.Table("tasks_performers").Select("user_id").Where("task_id = ?", taskID).Scan(&performerIDs)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	exist = false
	for _, id := range performerIDs {
		if id == req.PerformerID {
			exist = true
			break
		}
	}
	if !exist {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusNotAcceptable, messages.NotInTask)
	}
	result = configs.DB.Unscoped().Where("task_id = ? and user_id = ?", taskID, req.PerformerID).Delete(&models.TasksPerformers{})
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	var email string
	result = configs.DB.Table("users").Select("email").Where("id = ?", req.PerformerID).Scan(&email)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var projectTitle string
	result = configs.DB.Table("projects").Select("title").Where("id = ?", projectID).Scan(&projectTitle)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	var username string
	result = configs.DB.Table("users").Select("username").Where("id = ?", userID).Scan(&username)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	var taskTitle string
	result = configs.DB.Table("tasks").Select("title").Where("id = ?", taskID).Scan(&taskTitle)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	go func() {
		services.SendMail("PlanVerse Notification", fmt.Sprintf("you've been removed from %s task in %s project by %s!", taskTitle, projectTitle, username), []string{email})
	}()
	var task models.Task
	result = configs.DB.Where("id = ?", taskID).Preload("Performers").Find(&task)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	wsMessage := models.WSMessage{
		Type:    "update-task",
		Payload: task,
	}
	for conn, client := range models.Clients {
		for i := range members {
			if client == members[i] {
				wsErr := conn.WriteJSON(wsMessage)
				if wsErr != nil {
					delete(models.Clients, conn)
				}
			}
		}
	}
	models.SuccessRequests.WithLabelValues(method, endpoint).Inc()
	models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
	return ctx.JSON(http.StatusOK, messages.PerformerRemoved)
}

func EditTaskHandler(ctx echo.Context) error {
	startTime := time.Now()
	method := ctx.Request().Method
	endpoint := ctx.Request().URL.Path
	req := new(models.EditTaskRequest)
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
	taskID, err := strconv.Atoi(ctx.Param("task-id"))
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongTaskID)
	}
	var stateID int
	result := configs.DB.Table("tasks").Select("state_id").Where("id = ?", taskID).Scan(&stateID)
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	if stateID == 0 {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongTaskID)
	}
	var stateIDs []int
	result = configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
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
		return ctx.JSON(http.StatusNotAcceptable, messages.TaskNotInProject)
	}
	var task models.Task
	result = configs.DB.Where("id = ?", taskID).Find(&task)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	estimatedTime := 0
	deadline := time.Time{}
	year := 0
	month := 0
	day := 0
	if req.EstimatedTime == 0 {
		estimatedTime = task.EstimatedTime
	} else {
		estimatedTime = req.EstimatedTime
	}
	if req.Deadline == "" {
		deadline = task.Deadline
	} else {
		year, _ = strconv.Atoi(req.Deadline[:4])
		month, _ = strconv.Atoi(req.Deadline[5:7])
		day, _ = strconv.Atoi(req.Deadline[8:])
		deadline = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	}
	if req.Priority < 1 || req.Priority > 5 {
		req.Priority = task.Priority
	}
	task.Index = req.Index
	task.Title = req.Title
	task.BackGroundColor = req.BackGroundColor
	task.Description = req.Description
	task.EstimatedTime = estimatedTime
	task.Deadline = deadline
	task.Priority = req.Priority
	result = configs.DB.Save(&task)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	var members []int
	result = configs.DB.Table("projects_members").Select("user_id").Where("project_id = ?", projectID).Scan(&members)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	result = configs.DB.Where("id = ?", taskID).Preload("Performers").Find(&task)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	wsMessage := models.WSMessage{
		Type:    "update-task",
		Payload: task,
	}
	for conn, client := range models.Clients {
		for i := range members {
			if client == members[i] {
				wsErr := conn.WriteJSON(wsMessage)
				if wsErr != nil {
					delete(models.Clients, conn)
				}
			}
		}
	}
	models.SuccessRequests.WithLabelValues(method, endpoint).Inc()
	models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
	return ctx.JSON(http.StatusOK, messages.TaskEdited)
}

func DeleteTaskHandler(ctx echo.Context) error {
	startTime := time.Now()
	method := ctx.Request().Method
	endpoint := ctx.Request().URL.Path
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	taskID, err := strconv.Atoi(ctx.Param("task-id"))
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongTaskID)
	}
	var stateID int
	result := configs.DB.Table("tasks").Select("state_id").Where("id = ?", taskID).Scan(&stateID)
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	if stateID == 0 {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongTaskID)
	}
	var stateIDs []int
	result = configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
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
		return ctx.JSON(http.StatusNotAcceptable, messages.TaskNotInProject)
	}
	result = configs.DB.Unscoped().Where("task_id = ?", taskID).Delete(&models.TasksPerformers{})
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	result = configs.DB.Unscoped().Where("id = ?", taskID).Delete(&models.Task{})
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	var members []int
	result = configs.DB.Table("projects_members").Select("user_id").Where("project_id = ?", projectID).Scan(&members)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	var task models.Task
	result = configs.DB.Where("id = ?", taskID).Preload("Performers").Find(&task)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	wsMessage := models.WSMessage{
		Type:    "update-task",
		Payload: task,
	}
	for conn, client := range models.Clients {
		for i := range members {
			if client == members[i] {
				wsErr := conn.WriteJSON(wsMessage)
				if wsErr != nil {
					delete(models.Clients, conn)
				}
			}
		}
	}
	models.SuccessRequests.WithLabelValues(method, endpoint).Inc()
	models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
	return ctx.JSON(http.StatusOK, messages.TaskDeleted)
}

func GetTaskHandler(ctx echo.Context) error {
	startTime := time.Now()
	method := ctx.Request().Method
	endpoint := ctx.Request().URL.Path
	res := new(models.GetTaskResponse)
	taskID, err := strconv.Atoi(ctx.Param("task-id"))
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongTaskID)
	}
	var stateID int
	result := configs.DB.Table("tasks").Select("state_id").Where("id = ?", taskID).Scan(&stateID)
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	if stateID == 0 {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusBadRequest, messages.WrongTaskID)
	}
	var projectID int
	result = configs.DB.Table("states").Select("project_id").Where("id = ?", stateID).Scan(&projectID)
	if err != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
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
	result = configs.DB.Table("tasks").Select([]string{"index", "title", "back_ground_color", "description", "deadline", "estimated_time", "actual_time", "priority"}).Where("id = ?", taskID).Scan(res)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	var performerIDs []int
	result = configs.DB.Table("tasks_performers").Select("user_id").Where("task_id = ?", taskID).Scan(&performerIDs)
	if result.Error != nil {
		models.FailedRequests.WithLabelValues(method, endpoint).Inc()
		models.FailedDBRequests.WithLabelValues(method, endpoint).Inc()
		models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.SuccessDBRequests.WithLabelValues(method, endpoint).Inc()
	res.ID = taskID
	res.Performers = performerIDs
	models.SuccessRequests.WithLabelValues(method, endpoint).Inc()
	models.ResponseTime.WithLabelValues(method, endpoint).Observe(time.Since(startTime).Seconds())
	return ctx.JSON(http.StatusOK, res)
}

func CreateWSConnection(ctx echo.Context) error {
	log.Println("entered")
	userID, err := strconv.Atoi(ctx.Param("user-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	var upgrade = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrade.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	models.Clients[conn] = userID
	return nil
}
