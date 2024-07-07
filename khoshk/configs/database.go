package configs

import (
	"PlanVerse/models"
	"fmt"
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	onceDB sync.Once
)

func ConnectToDatabase() {
	onceDB.Do(func() {
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN: fmt.Sprintf(
				"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
				os.Getenv("POSTGRES_HOST"),
				os.Getenv("POSTGRES_PORT"),
				os.Getenv("POSTGRES_USER"),
				os.Getenv("POSTGRES_PASSWORD"),
				os.Getenv("POSTGRES_DB"),
			),
		}), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		err = db.SetupJoinTable(&models.User{}, "Projects", &models.ProjectsMembers{})
		if err != nil {
			log.Fatalf("Failed to set up join table: %v", err)
		}
		err = db.SetupJoinTable(&models.Project{}, "Members", &models.ProjectsMembers{})
		if err != nil {
			log.Fatalf("Failed to set up join table: %v", err)
		}
		err = db.SetupJoinTable(&models.Project{}, "InvitedMembers", &models.InvitedMembers{})
		if err != nil {
			log.Fatalf("Failed to set up join table: %v", err)
		}
		err = db.SetupJoinTable(&models.Task{}, "Performers", &models.TasksPerformers{})
		if err != nil {
			log.Fatalf("Failed to set up join table: %v", err)
		}
		err = db.AutoMigrate(&models.User{})
		if err != nil {
			log.Fatalf("Failed to migrate user tabel: %v", err)
		}
		err = db.AutoMigrate(&models.Project{})
		if err != nil {
			log.Fatalf("Failed to migrate project tabel: %v", err)
		}
		err = db.AutoMigrate(&models.State{})
		if err != nil {
			log.Fatalf("Failed to migrate state tabel: %v", err)
		}
		err = db.AutoMigrate(&models.Task{})
		if err != nil {
			log.Fatalf("Failed to migrate task tabel: %v", err)
		}
		err = db.AutoMigrate(&models.Label{})
		if err != nil {
			log.Fatalf("Failed to migrate label tabel: %v", err)
		}
		err = db.AutoMigrate(&models.Comment{})
		if err != nil {
			log.Fatalf("Failed to migrate comment tabel: %v", err)
		}
		err = db.AutoMigrate(&models.JoinLink{})
		if err != nil {
			log.Fatalf("Failed to migrate join_link tabel: %v", err)
		}
		DB = db
	})
}
