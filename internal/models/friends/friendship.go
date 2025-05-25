package friends

import (
	"time"

	"gorm.io/gorm"
)

type Friendship struct {
    FriendshipID        uint64    `json:"friendship_id" gorm:"primaryKey;autoIncrement"`
    User1ID   string    `json:"user1_id" gorm:"type:char(36);not null;index:idx_friendships_user1;uniqueIndex:unique_friendship;constraint:OnDelete:CASCADE"`
    User2ID   string    `json:"user2_id" gorm:"type:char(36);not null;index:idx_friendships_user2;uniqueIndex:unique_friendship;constraint:OnDelete:CASCADE"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (Friendship) TableName() string {
    return "friendships"
}

func (f *Friendship) BeforeCreate(tx *gorm.DB) (err error) {
	// Ensure User1ID is always less than User2ID to maintain uniqueness
	if f.User1ID > f.User2ID {
		f.User1ID, f.User2ID = f.User2ID, f.User1ID
	}
	return nil
}