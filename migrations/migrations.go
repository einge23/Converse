package migrations

import (
	"converse/internal/db"
	"converse/internal/models"
	"converse/internal/models/friends"
	"log"
)

func RunMigrations() error {
    database := db.GetDB()

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