package migrations

import (
	"converse/internal/db"
	"converse/internal/models"
	"converse/internal/models/friends"
	"log"
)

func RunMigrations() error {
    database := db.GetDB()

    database.Exec("SET FOREIGN_KEY_CHECKS = 0")
    database.Exec("ALTER TABLE sessions DROP FOREIGN KEY IF EXISTS fk_sessions_user")
    database.Exec("ALTER TABLE users DROP FOREIGN KEY IF EXISTS fk_sessions_user")
    database.Exec("SET FOREIGN_KEY_CHECKS = 1")

    // Run schema migrations
    err := database.AutoMigrate(
        &models.User{},
        &models.Session{},
        &friends.FriendRequest{},
        &friends.Friendship{},
        &models.DirectMessageThread{},
        &models.Message{},
        // &models.Room{},
        // &models.RoomMember{},
    )
    if err != nil {
        return err
    }

    log.Println("Migrations completed successfully")
    return nil
}