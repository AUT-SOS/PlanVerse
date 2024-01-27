package controllers

import (
	"PlanVerse/configs"
	"PlanVerse/helpers"
	"PlanVerse/messages"
	"PlanVerse/models"
	"PlanVerse/services"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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

func ShareProjectHandler(ctx echo.Context) error {
	req := new(models.ShareProjectRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	var user models.User
	result := configs.DB.Select("username").Where("id = ?", userID).Find(&user)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var joinLink models.JoinLink
	result = configs.DB.Select("link").Where("project_id = ?", req.ProjectID).Find(&joinLink)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var project models.Project
	result = configs.DB.Where("id = ?", req.ProjectID).Preload("InvitedMembers").Find(&project)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var projectMembers []helpers.ProjectMember
	result = configs.DB.Table("projects").Select([]string{"users.email"}).Joins("inner join projects_members on projects.id = projects_members.project_id").Joins("inner join users on users.id = projects_members.user_id").Where("projects.id = ?", req.ProjectID).Scan(&projectMembers)
	for i, _ := range req.Emails {
		for j, _ := range projectMembers {
			if projectMembers[j].Email == req.Emails[i] {
				return ctx.JSON(http.StatusNotAcceptable, messages.AlreadyMember)
			}
		}
	}
	for i, _ := range req.Emails {
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
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	var projectID helpers.ProjectID
	result := configs.DB.Table("join_links").Select("project_id").Where("link = ?", req.Link).Scan(&projectID)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	var project models.Project
	result = configs.DB.Where("id = ?", projectID.ProjectID).Preload("Members").Find(&project)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	res.ProjectID = projectID.ProjectID
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
	var projectIDStructs []helpers.ProjectID
	result = configs.DB.Table("invited_members").Select("project_id").Where("user_id = ?", userID).Find(&projectIDStructs)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	isInvited := false
	for i := 0; i < len(projectIDStructs); i++ {
		if projectIDStructs[i].ProjectID == projectID {
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
	result = configs.DB.Table("projects_members").Where("project_id = ? and user_id = ?", projectID, userID).Update("is_admin", false)
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
	var showRole helpers.ShowRole
	result := configs.DB.Table("projects_members").Where("project_id = ? and user_id = ?", projectID, memberID).Scan(&showRole)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
	}
	if showRole.IsAdmin {
		return ctx.JSON(http.StatusNotAcceptable, messages.AlreadyAdmin)
	}
	result = configs.DB.Table("projects_members").Where("project_id = ? and user_id = ?", projectID, memberID).Update("is_admin", true)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
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
	var owner helpers.Owner
	result := configs.DB.Table("projects").Select("owner_id").Where("id = ?", projectID).Scan(&owner)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.WrongProjectID)
	}
	if userID != owner.OwnerID {
		return ctx.JSON(http.StatusNotAcceptable, messages.OwnerAccess)
	}
	if adminID == owner.OwnerID {
		return ctx.JSON(http.StatusNotAcceptable, messages.OwnerChange)
	}
	var showRole helpers.ShowRole
	result = configs.DB.Table("projects_members").Where("project_id = ? and user_id = ?", projectID, adminID).Scan(&showRole)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
	}
	if !showRole.IsAdmin {

	}
	result = configs.DB.Table("projects_members").Where("project_id = ? and user_id = ?", projectID, adminID).Update("is_admin", false)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.AlreadyMemberRole)
	}
	return ctx.JSON(http.StatusOK, messages.AdminRoleChanged)
}
