package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title           string    `gorm:"not null"`
	Description     string    `gorm:"not null"`
	BackGroundColor string    `gorm:"not null"`
	StateID         int       `gorm:"not null"`
	State           State     `gorm:"foreignKey:StateID"`
	Performers      []User    `gorm:"many2many:tasks_performers"`
	Comments        []Comment `gorm:"foreignKey:TaskID"`
	Labels          []Label   `gorm:"many2many:tasks_labels"`
}

type TaskShow struct {
	ID              int
	Title           string
	BackGroundColor string
	Performers      []int
}

type CreateTaskRequest struct {
	StateID         int    `json:"state_id"`
	Title           string `json:"title"`
	BackGroundColor string `json:"back_ground_color"`
	Description     string `json:"description"`
}

type ChangeTaskStateRequest struct {
	TaskID  int `json:"task_id"`
	StateID int `json:"state_id"`
}

type PerformerRequest struct {
	TaskID      int `json:"task_id"`
	PerformerID int `json:"performer_id"`
}

type EditTaskRequest struct {
	TaskID          int    `json:"task_id"`
	Title           string `json:"title"`
	BackGroundColor string `json:"back_ground_color"`
	Description     string `json:"description"`
}

type CreateTaskResponse struct {
	TaskID int `json:"task_id"`
}
