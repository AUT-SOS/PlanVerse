package models

type Task struct {
	ID              uint      `gorm:"primaryKey"`
	Title           string    `gorm:"not null"`
	Description     string    `gorm:"not null"`
	BackGroundColor string    `gorm:"not null"`
	StateID         int       `gorm:"not null"`
	State           State     `gorm:"foreignKey:StateID"`
	Performers      []User    `gorm:"many2many:tasks_performers"`
	Comments        []Comment `gorm:"foreignKey:TaskID"`
	Labels          []Label   `gorm:"many2many:tasks_labels"`
}
