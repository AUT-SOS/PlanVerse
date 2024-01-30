package controllers

import (
	"PlanVerse/configs"
	"PlanVerse/helpers"
	"PlanVerse/messages"
	"PlanVerse/models"
	"PlanVerse/services"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"sync"
	"time"
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
	defaultStates := make([]models.State, 3)
	defaultStates[0] = models.State{
		Title:           "To-Do",
		BackGroundColor: "D0D613",
		AdminAccess:     false,
	}
	defaultStates[1] = models.State{
		Title:           "Doing",
		BackGroundColor: "DE731A",
		AdminAccess:     false,
	}
	defaultStates[2] = models.State{
		Title:           "Done",
		BackGroundColor: "54D826",
		AdminAccess:     true,
	}
	newProject := models.Project{
		Title:         req.Title,
		Description:   req.Description,
		BackGroundPic: req.Picture,
		OwnerID:       userID,
		MembersNumber: 1,
		Members:       []models.User{user},
		JoinLink:      joinLink,
		States:        defaultStates,
	}
	if err = configs.DB.Create(&newProject).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToCreateProject)
	}
	result = configs.DB.Table("projects_members").Where("project_id = ?", newProject.ID).Update("is_admin", true)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	result = configs.DB.Table("projects_members").Where("project_id = ?", newProject.ID).Update("promotion_time", time.Now())
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	res.ProjectID = int(newProject.ID)
	res.Link = joinLink.Link
	return ctx.JSON(http.StatusOK, res)
}

func ShareProjectHandler(ctx echo.Context) error {
	req := new(models.ShareProjectRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	var projectIDs []int
	result := configs.DB.Table("projects_members").Select("project_id").Where("user_id = ?", userID).Scan(&projectIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if len(projectIDs) == 0 {
		return ctx.JSON(http.StatusNotAcceptable, messages.UserNoProject)
	}
	exist := false
	for _, id := range projectIDs {
		if id == projectID {
			exist = true
			break
		}
	}
	if !exist {
		return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
	}
	var user models.User
	result = configs.DB.Select("username").Where("id = ?", userID).Find(&user)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var joinLink models.JoinLink
	result = configs.DB.Select("link").Where("project_id = ?", projectID).Find(&joinLink)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var project models.Project
	result = configs.DB.Where("id = ?", projectID).Preload("InvitedMembers").Find(&project)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var projectMembersEmail []string
	result = configs.DB.Table("projects").Select([]string{"users.email"}).Joins("inner join projects_members on projects.id = projects_members.project_id").Joins("inner join users on users.id = projects_members.user_id").Where("projects.id = ?", projectID).Scan(&projectMembersEmail)
	for i := range req.Emails {
		for j := range projectMembersEmail {
			if projectMembersEmail[j] == req.Emails[i] {
				return ctx.JSON(http.StatusNotAcceptable, messages.AlreadyMember)
			}
		}
	}
	for i := range req.Emails {
		go func(index int) {
			services.SendMail("PlanVerse Invitation", fmt.Sprintf("you've been invited to %s project by %s!\nclick the link below to join to project:\n%s", project.Title, user.Username, joinLink.Link), []string{req.Emails[index]})
		}(i)
	}
	var users []models.User
	result = configs.DB.Where("email in ?", req.Emails).Find(&users)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	for _, u := range users {
		project.InvitedMembers = append(project.InvitedMembers, u)
	}
	result = configs.DB.Save(&project)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	return ctx.JSON(http.StatusOK, messages.SentInvitationEmail)
}

func ShowProjectHandler(ctx echo.Context) error {
	req := new(models.ShowProjectRequest)
	res := new(models.ShowProjectResponse)
	if err := json.NewDecoder(ctx.Request().Body).Decode(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	var projectID int
	result := configs.DB.Table("join_links").Select("project_id").Where("link = ?", req.Link).Scan(&projectID)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if projectID == 0 {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongLink)
	}
	var project models.Project
	result = configs.DB.Where("id = ?", projectID).Preload("Members").Find(&project)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	res.ProjectID = projectID
	res.Title = project.Title
	res.BackGroundPic = project.BackGroundPic
	res.MembersNumber = project.MembersNumber
	for i := 0; i < len(project.Members); i++ {
		if i == 3 {
			break
		}
		user := models.MemberInfo{
			Username:   project.Members[i].Username,
			ProfilePic: project.Members[i].ProfilePic,
		}
		res.Members = append(res.Members, user)
	}
	return ctx.JSON(http.StatusOK, res)
}

func JoinProjectHandler(ctx echo.Context) error {
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	var project models.Project
	result := configs.DB.Where("id = ?", projectID).Preload("Members").Find(&project)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if project.ID == 0 {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	for _, member := range project.Members {
		if int(member.ID) == userID {
			return ctx.JSON(http.StatusNotAcceptable, messages.AlreadyJoined)
		}
	}
	var user models.User
	result = configs.DB.Where("id = ?", userID).Find(&user)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var projectIDs []int
	result = configs.DB.Table("invited_members").Select("project_id").Where("user_id = ?", userID).Find(&projectIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	isInvited := false
	for i := 0; i < len(projectIDs); i++ {
		if projectIDs[i] == projectID {
			isInvited = true
			break
		}
	}
	if !isInvited {
		return ctx.JSON(http.StatusNotAcceptable, messages.Uninvited)
	}
	result = configs.DB.Unscoped().Where("project_id = ? and user_id = ?", projectID, userID).Delete(&models.InvitedMembers{})
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	project.Members = append(project.Members, user)
	project.MembersNumber += 1
	result = configs.DB.Save(&project)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	return ctx.JSON(http.StatusOK, messages.UserAddedToProject)
}

func ChangeRoleMemberHandler(ctx echo.Context) error {
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	memberID, err := strconv.Atoi(ctx.Param("user-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongUserID)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	var projectIDs []int
	result := configs.DB.Table("projects_members").Select("project_id").Where("user_id = ?", memberID).Scan(&projectIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if len(projectIDs) == 0 {
		return ctx.JSON(http.StatusNotAcceptable, messages.UserNoProject)
	}
	exist := false
	for _, id := range projectIDs {
		if id == projectID {
			exist = true
			break
		}
	}
	if !exist {
		return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
	}
	var isAdmin bool
	result = configs.DB.Table("projects_members").Select("is_admin").Where("project_id = ? and user_id = ?", projectID, memberID).Scan(&isAdmin)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if isAdmin {
		return ctx.JSON(http.StatusNotAcceptable, messages.AlreadyAdmin)
	}
	var email string
	result = configs.DB.Table("users").Select("email").Where("id = ?", memberID).Scan(&email)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var projectTitle string
	result = configs.DB.Table("projects").Select("title").Where("id = ?", projectID).Scan(&projectTitle)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var username string
	result = configs.DB.Table("users").Select("username").Where("id = ?", userID).Scan(&username)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	result = configs.DB.Table("projects_members").Where("project_id = ? and user_id = ?", projectID, memberID).Update("is_admin", true)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	result = configs.DB.Table("projects_members").Where("project_id = ? and user_id = ?", projectID, memberID).Update("promotion_time", time.Now())
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	go func() {
		services.SendMail("PlanVerse Notification", fmt.Sprintf("you've been promoted to admin role in %s project by %s!", projectTitle, username), []string{email})
	}()
	return ctx.JSON(http.StatusOK, messages.MemberRoleChanged)
}

func ChangeRoleAdminHandler(ctx echo.Context) error {
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	adminID, err := strconv.Atoi(ctx.Param("user-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongUserID)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	var projectIDs []int
	result := configs.DB.Table("projects_members").Select("project_id").Where("user_id = ?", adminID).Scan(&projectIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if len(projectIDs) == 0 {
		return ctx.JSON(http.StatusNotAcceptable, messages.UserNoProject)
	}
	exist := false
	for _, id := range projectIDs {
		if id == projectID {
			exist = true
			break
		}
	}
	if !exist {
		return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
	}
	var ownerID int
	result = configs.DB.Table("projects").Select("owner_id").Where("id = ?", projectID).Scan(&ownerID)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if ownerID == 0 {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	if userID != ownerID {
		return ctx.JSON(http.StatusNotAcceptable, messages.OwnerAccess)
	}
	if adminID == ownerID {
		return ctx.JSON(http.StatusNotAcceptable, messages.OwnerChange)
	}
	var isAdmin bool
	result = configs.DB.Table("projects_members").Select("is_admin").Where("project_id = ? and user_id = ?", projectID, adminID).Scan(&isAdmin)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if !isAdmin {
		return ctx.JSON(http.StatusNotAcceptable, messages.AlreadyMemberRole)
	}
	var email string
	result = configs.DB.Table("users").Select("email").Where("id = ?", adminID).Scan(&email)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var projectTitle string
	result = configs.DB.Table("projects").Select("title").Where("id = ?", projectID).Scan(&projectTitle)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var username string
	result = configs.DB.Table("users").Select("username").Where("id = ?", userID).Scan(&username)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	result = configs.DB.Table("projects_members").Where("project_id = ? and user_id = ?", projectID, adminID).Update("is_admin", false)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	result = configs.DB.Table("projects_members").Where("project_id = ? and user_id = ?", projectID, adminID).Update("promotion_time", time.Time{})
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	go func() {
		services.SendMail("PlanVerse Notification", fmt.Sprintf("you've been demoted to member role in %s project by %s!", projectTitle, username), []string{email})
	}()
	return ctx.JSON(http.StatusOK, messages.AdminRoleChanged)
}

func GetProjectHandler(ctx echo.Context) error {
	res := new(models.GetProjectResponse)
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	var projectIDs []int
	result := configs.DB.Table("projects_members").Select([]string{"project_id"}).Where("user_id = ?", userID).Scan(&projectIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if len(projectIDs) == 0 {
		return ctx.JSON(http.StatusNotAcceptable, messages.UserNoProject)
	}
	isMember := false
	for _, project := range projectIDs {
		if project == projectID {
			isMember = true
			break
		}
	}
	if !isMember {
		return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
	}
	result = configs.DB.Table("projects").Select([]string{"title", "back_ground_pic", "description", "owner_id", "members_number"}).Where("id = ?", projectID).Scan(res)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var joinLink models.JoinLink
	result = configs.DB.Table("join_links").Select("link").Where("project_id = ?", projectID).Scan(&joinLink)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	res.ID = projectID
	res.Link = joinLink.Link
	return ctx.JSON(http.StatusOK, res)
}

func GetProjectMembersHandler(ctx echo.Context) error {
	var res []models.GetMemberResponse
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	var projectIDs []int
	result := configs.DB.Table("projects_members").Select("project_id").Where("user_id = ?", userID).Scan(&projectIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if len(projectIDs) == 0 {
		return ctx.JSON(http.StatusNotAcceptable, messages.UserNoProject)
	}
	isMember := false
	for _, id := range projectIDs {
		if id == projectID {
			isMember = true
			break
		}
	}
	if !isMember {
		return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
	}
	var project models.Project
	result = configs.DB.Where("id = ?", projectID).Preload("Members").Find(&project)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	res = make([]models.GetMemberResponse, len(project.Members))
	showRoles := make([]bool, len(project.Members))
	var wg sync.WaitGroup
	for i := range project.Members {
		wg.Add(1)
		go func(index int, wg *sync.WaitGroup) {
			defer wg.Done()
			configs.DB.Table("projects_members").Select("is_admin").Where("project_id = ? and user_id = ?", projectID, project.Members[index].ID).Scan(&showRoles[index])
			res[index] = models.GetMemberResponse{
				ID:         int(project.Members[index].ID),
				Username:   project.Members[index].Username,
				Email:      project.Members[index].Email,
				ProfilePic: project.Members[index].ProfilePic,
				IsAdmin:    showRoles[index],
			}
		}(i, &wg)
	}
	wg.Wait()
	return ctx.JSON(http.StatusOK, res)
}

func LeaveProjectHandler(ctx echo.Context) error {
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	var projectIDs []int
	result := configs.DB.Table("projects_members").Select("project_id").Where("user_id = ?", userID).Scan(&projectIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if len(projectIDs) == 0 {
		return ctx.JSON(http.StatusNotAcceptable, messages.UserNoProject)
	}
	isMember := false
	for _, id := range projectIDs {
		if id == projectID {
			isMember = true
			break
		}
	}
	if !isMember {
		return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
	}
	var membersNumber int
	result = configs.DB.Table("projects").Select("members_number").Where("id = ?", projectID).Scan(&membersNumber)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if membersNumber == 0 {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	result = configs.DB.Unscoped().Where("user_id = ? and project_id = ?", userID, projectID).Delete(&models.ProjectsMembers{})
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	membersNumber -= 1
	if membersNumber == 0 {
		result = configs.DB.Unscoped().Where("project_id = ?", projectID).Delete(&models.InvitedMembers{})
		if result.Error != nil {
			return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
		}
		var stateIDs []int
		result = configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
		if result.Error != nil {
			return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
		}
		var wg sync.WaitGroup
		for i := range stateIDs {
			wg.Add(1)
			go func(index int, wg *sync.WaitGroup) {
				defer wg.Done()
				var taskIDs []int
				result = configs.DB.Table("tasks").Select("id").Where("state_id = ?", stateIDs[index]).Scan(&taskIDs)
				for j := range taskIDs {
					configs.DB.Unscoped().Where("task_id = ?", taskIDs[j]).Delete(&models.TasksPerformers{})
				}
				result = configs.DB.Unscoped().Where("state_id = ?", stateIDs[index]).Delete(&models.Task{})
			}(i, &wg)
		}
		wg.Wait()
		result = configs.DB.Unscoped().Where("project_id = ?", projectID).Delete(&models.State{})
		if result.Error != nil {
			return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
		}
		result = configs.DB.Unscoped().Where("project_id = ?", projectID).Delete(&models.JoinLink{})
		if result.Error != nil {
			return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
		}
		result = configs.DB.Unscoped().Where("id = ?", projectID).Delete(&models.Project{})
		if result.Error != nil {
			return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
		}
		return ctx.JSON(http.StatusOK, messages.UserLeavedProject)
	}
	var ownerID int
	result = configs.DB.Table("projects").Select("owner_id").Where("id = ?", projectID).Scan(&ownerID)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if ownerID == userID {
		newOwnerID, err := helpers.DetectMin(projectID)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
		}
		if newOwnerID != 0 {
			result = configs.DB.Table("projects").Where("id = ?", projectID).Update("owner_id", newOwnerID)
			if result.Error != nil {
				return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
			}
		}
		result = configs.DB.Table("projects").Where("id = ?", projectID).Update("members_number", membersNumber)
		if result.Error != nil {
			return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
		}
		return ctx.JSON(http.StatusOK, messages.UserLeavedProject)
	}
	result = configs.DB.Table("projects").Where("id = ?", projectID).Update("members_number", membersNumber)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	return ctx.JSON(http.StatusOK, messages.UserLeavedProject)
}

func EditProjectHandler(ctx echo.Context) error {
	req := new(models.EditProjectRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	var project models.Project
	result := configs.DB.Where("id = ?", projectID).Find(&project)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	project.Title = req.Title
	project.BackGroundPic = req.BackGroundPic
	project.Description = req.Description
	result = configs.DB.Save(&project)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	return ctx.JSON(http.StatusOK, messages.ProjectEdited)
}

func DeleteProjectHandler(ctx echo.Context) error {
	projectID, err := strconv.Atoi(ctx.Param("project-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	var ownerID int
	result := configs.DB.Table("projects").Select("owner_id").Where("id = ?", projectID).Scan(&ownerID)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if ownerID == 0 {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	if userID != ownerID {
		return ctx.JSON(http.StatusNotAcceptable, messages.OwnerAccess)
	}
	result = configs.DB.Unscoped().Where("project_id = ?", projectID).Delete(&models.InvitedMembers{})
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	result = configs.DB.Unscoped().Where("project_id = ?", projectID).Delete(&models.ProjectsMembers{})
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var stateIDs []int
	result = configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var wg sync.WaitGroup
	for i := range stateIDs {
		wg.Add(1)
		go func(index int, wg *sync.WaitGroup) {
			defer wg.Done()
			var taskIDs []int
			configs.DB.Table("tasks").Select("id").Where("state_id = ?", stateIDs[index]).Scan(&taskIDs)
			for j := range taskIDs {
				configs.DB.Unscoped().Where("task_id = ?", taskIDs[j]).Delete(&models.TasksPerformers{})
			}
			configs.DB.Unscoped().Where("state_id = ?", stateIDs[index]).Delete(&models.Task{})
		}(i, &wg)
	}
	wg.Wait()
	result = configs.DB.Unscoped().Where("project_id = ?", projectID).Delete(&models.State{})
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	result = configs.DB.Unscoped().Where("project_id = ?", projectID).Delete(&models.JoinLink{})
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	result = configs.DB.Unscoped().Where("id = ?", projectID).Delete(&models.Project{})
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	return ctx.JSON(http.StatusOK, messages.ProjectDeleted)
}
