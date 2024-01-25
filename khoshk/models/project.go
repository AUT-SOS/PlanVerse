package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Title         string   `gorm:"not null"`
	Description   string   `gorm:"not null"`
	BackGroundPic string   `gorm:"not null"`
	OwnerID       int      `gorm:"not null"`
	Members       []User   `gorm:"many2many:projects_members"`
	MembersNumber int      `gorm:"not null"`
	States        []State  `gorm:"foreignKey:ProjectID"`
	JoinLink      JoinLink `gorm:"foreignKey:ProjectID"`
}

type ResponseList struct {
	ProjectID     int    `json:"project_id"`
	Title         string `json:"project_title"`
	Picture       string `json:"project_picture"`
	MembersNumber int    `json:"members_number"`
}
