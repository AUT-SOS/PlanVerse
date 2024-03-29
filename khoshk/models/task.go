package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title           string    `gorm:"not null"`
	Description     string    `gorm:"not null"`
	BackGroundColor string    `gorm:"not null"`
	Index           int       `gorm:"not null"`
	StateID         int       `gorm:"not null"`
	State           State     `gorm:"foreignKey:StateID"`
	Performers      []User    `gorm:"many2many:tasks_performers"`
	Comments        []Comment `gorm:"foreignKey:TaskID"`
	Labels          []Label   `gorm:"many2many:tasks_labels"`
}

type TaskShow struct {
	ID              int    `json:"task_id"`
	Index           int    `json:"index"`
	Title           string `json:"title"`
	BackGroundColor string `json:"back_ground_color"`
	Performers      []int  `json:"performers"`
	Description     string `json:"description"`
}

type CreateTaskRequest struct {
	StateID         int    `json:"state_id"`
	Index           int    `json:"index"`
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
	Index           int    `json:"index"`
	Title           string `json:"title"`
	BackGroundColor string `json:"back_ground_color"`
	Description     string `json:"description"`
}

type CreateTaskResponse struct {
	TaskID int `json:"task_id"`
}

type GetTaskResponse struct {
	ID              int    `json:"task_id"`
	Index           int    `json:"index"`
	Title           string `json:"title"`
	BackGroundColor string `json:"back_ground_color"`
	Description     string `json:"description"`
	Performers      []int  `json:"performers"`
}
