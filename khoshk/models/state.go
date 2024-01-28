package models

import "gorm.io/gorm"

type State struct {
	gorm.Model
	Title           string  `gorm:"not null"`
	BackGroundColor string  `gorm:"not null"`
	AdminAccess     bool    `gorm:"not null"`
	ProjectID       int     `gorm:"not null"`
	Project         Project `gorm:"foreignKey:ProjectID"`
	Tasks           []Task  `gorm:"foreignKey:StateID"`
}

type CreateStateRequest struct {
	Title           string `json:"title"`
	BackGroundColor string `json:"back_ground_color"`
	AdminAccess     bool   `json:"admin_access"`
}

type StateObject struct {
	ID              int    `json:"state_id"`
	Title           string `json:"title"`
	BackGroundColor string `json:"back_ground_color"`
	AdminAccess     bool   `json:"admin_access"`
}

type CreateStateResponse struct {
	StateID int `json:"state_id"`
}

type StateListResponse struct {
	ID              int        `json:"state_id"`
	Title           string     `json:"title"`
	BackGroundColor string     `json:"back_ground_color"`
	AdminAccess     bool       `json:"admin_access"`
	Tasks           []TaskShow `json:"tasks"`
}

type GetStateResponse struct {
	ID              int    `json:"state_id"`
	ProjectID       int    `json:"project_id"`
	Title           string `json:"title"`
	BackGroundColor string `json:"back_ground_color"`
	AdminAccess     bool   `json:"admin_access"`
}
