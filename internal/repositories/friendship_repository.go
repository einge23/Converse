package repositories

import (
	"converse/internal/db"
	"converse/internal/models"

	"gorm.io/gorm"
)

type FriendRepository struct {
	db *gorm.DB
}

func NewFriendRepository() *FriendRepository {
	return &FriendRepository{
		db: db.GetDB(),
	}
}

func (r *FriendRepository) GetFriends(userID string) ([]*models.PublicUser, error) {
    var users []*models.PublicUser
    err := r.db.Table("friendships f").
        Select(`u.user_id, u.username, u.email, u.display_name, u.avatar_url, u.status, u.last_active_at, 
                u.created_at, u.updated_at, u.deleted_at, dm.thread_id as dm_thread_id`).
        Joins(`JOIN users u ON (
            CASE 
                WHEN f.user1_id = ? THEN u.user_id = f.user2_id
                WHEN f.user2_id = ? THEN u.user_id = f.user1_id
            END
        )`, userID, userID).
        Joins(`LEFT JOIN direct_message_threads dm ON (
            (dm.user1_id = ? AND dm.user2_id = u.user_id) OR 
            (dm.user2_id = ? AND dm.user1_id = u.user_id)
        )`, userID, userID).
        Where("(f.user1_id = ? OR f.user2_id = ?) AND u.deleted_at IS NULL", userID, userID).
        Order("u.username ASC").
        Scan(&users).Error
    
    if err != nil {
        return nil, err
    }
    
    return users, nil
}