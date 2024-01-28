package models

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Title          string   `gorm:"not null"`
	Description    string   `gorm:"not null"`
	BackGroundPic  string   `gorm:"not null"`
	OwnerID        int      `gorm:"not null"`
	MembersNumber  int      `gorm:"not null"`
	InvitedMembers []User   `gorm:"many2many:invited_members"`
	Members        []User   `gorm:"many2many:projects_members"`
	States         []State  `gorm:"foreignKey:ProjectID"`
	JoinLink       JoinLink `gorm:"foreignKey:ProjectID"`
}

type CreateProjectRequest struct {
	Title       string `json:"title"`
	Picture     string `json:"picture"`
	Description string `json:"description"`
}

type ShareProjectRequest struct {
	Emails []string `json:"emails"`
}

type ShowProjectRequest struct {
	Link string `json:"join_link"`
}

type ProjectListResponse struct {
	ID            int    `json:"project_id"`
	Title         string `json:"title"`
	BackGroundPic string `json:"picture"`
	MembersNumber int    `json:"members_number"`
	IsAdmin       bool   `json:"is_admin"`
}

type EditProjectRequest struct {
	Title         string `json:"title"`
	BackGroundPic string `json:"picture"`
	Description   string `json:"description"`
}

type CreateProjectResponse struct {
	ProjectID int    `json:"project_id"`
	Link      string `json:"join_link"`
}

type ShowProjectResponse struct {
	ProjectID     int          `json:"project_id"`
	Title         string       `json:"title"`
	BackGroundPic string       `json:"picture"`
	MembersNumber int          `json:"members_number"`
	Members       []MemberInfo `json:"members"`
}

type GetProjectResponse struct {
	ID            int    `json:"project_id"`
	Title         string `json:"title"`
	BackGroundPic string `json:"picture"`
	Description   string `json:"description"`
	OwnerID       int    `json:"owner_id"`
	MembersNumber int    `json:"members_number"`
}

type GetMemberResponse struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	ProfilePic string `json:"profile_pic"`
	IsAdmin    bool   `json:"is_admin"`
}

type MemberInfo struct {
	Username   string
	ProfilePic string
}
