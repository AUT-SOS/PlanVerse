package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Title         string   `gorm:"not null"`
	Description   string   `gorm:"not null"`
	BackGroundPic string   `gorm:"not null"`
	OwnerID       int      `gorm:"not null"`
	MembersNumber int      `gorm:"not null"`
	Members       []User   `gorm:"many2many:projects_members"`
	States        []State  `gorm:"foreignKey:ProjectID"`
	JoinLink      JoinLink `gorm:"foreignKey:ProjectID"`
}

type CreateProjectRequest struct {
	Title       string `json:"title"`
	Picture     string `json:"picture"`
	Description string `json:"description"`
}

type ShareProjectRequest struct {
	ProjectID string   `json:"project_id"`
	Emails    []string `json:"emails"`
}

type ProjectListResponse struct {
	ID            int    `json:"project_id"`
	Title         string `json:"title"`
	BackGroundPic string `json:"picture"`
	MembersNumber int    `json:"members_number"`
	IsAdmin       bool   `json:"is_admin"`
}

type CreateProjectResponse struct {
	ProjectID int    `json:"project_id"`
	Link      string `json:"join_link"`
	Message   string `json:"message"`
}
