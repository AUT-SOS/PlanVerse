package models

type Comment struct {
	ID       uint   `gorm:"primaryKey"`
	Text     string `gorm:"not null"`
	SenderID int    `gorm:"not null"`
	TaskID   int    `gorm:"not null"`
	Task     Task   `gorm:"foreignKey:TaskID"`
}
