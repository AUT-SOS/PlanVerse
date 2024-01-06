package models

type JoinLink struct {
	ID        int    `gorm:"primaryKey"`
	Link      string `gorm:"unique;not null"`
	ProjectID int    `gorm:"notnull"`
}
