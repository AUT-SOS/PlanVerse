package models

type Label struct {
	ID            uint   `gorm:"primaryKey"`
	Title         string `gorm:"not null"`
	BackGroundPic string `gorm:"not null"`
	Tasks         []Task `gorm:"many2many:tasks_labels"`
}
