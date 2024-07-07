package models

import (
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model
	Title           string    `gorm:"not null"`
	Description     string    `gorm:"not null"`
	BackGroundColor string    `gorm:"not null"`
	EstimatedTime   int       `gorm:"not null"`
	ActualTime      int       `gorm:"not null"`
	Deadline        time.Time `gorm:"not null"`
	Priority        int       `gorm:"not null"`
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
	Deadline        string `json:"deadline"`
	EstimatedTime   int    `json:"estimated_time"`
	ActualTime      int    `json:"actual_time"`
	Priority        int    `json:"priority"`
}

type CreateTaskRequest struct {
	Index           int    `json:"index"`
	Title           string `json:"title"`
	BackGroundColor string `json:"back_ground_color"`
	Description     string `json:"description"`
	Deadline        string `json:"deadline"`
	EstimatedTime   int    `json:"estimated_time"`
	Priority        int    `json:"priority"`
}

type ChangeTaskStateRequest struct {
	StateID int `json:"state_id"`
}

type PerformerRequest struct {
	PerformerID int `json:"performer_id"`
}

type EditTaskRequest struct {
	Index           int    `json:"index"`
	Title           string `json:"title"`
	BackGroundColor string `json:"back_ground_color"`
	Description     string `json:"description"`
	Deadline        string `json:"deadline"`
	EstimatedTime   int    `json:"estimated_time"`
	Priority        int    `json:"priority"`
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
	Deadline        string `json:"deadline"`
	EstimatedTime   int    `json:"estimated_time"`
	ActualTime      int    `json:"actual_time"`
	Priority        int    `json:"priority"`
}

type WSMessage struct {
	Type    string `json:"type"`
	Payload Task   `json:"payload"`
}
