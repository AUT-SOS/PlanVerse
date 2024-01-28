package models

type TasksPerformers struct {
	UserID uint `gorm:"primaryKey"`
	TaskID uint `gorm:"primaryKey"`
}
