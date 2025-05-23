package migrations

import (
	"converse/internal/db"
	"converse/internal/models"
	"log"
)

func RunMigrations() error {
    database := db.GetDB()

    err := database.AutoMigrate(&models.User{}, &models.Session{})
    if err != nil {
        return err
    }

    log.Println("Migrations completed successfully")
    return nil
}