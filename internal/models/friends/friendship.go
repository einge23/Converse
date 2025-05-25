package friends

import "time"

type Friendship struct {
    ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
    User1ID   string    `json:"user1_id" gorm:"type:char(36);not null;index:idx_friendships_user1;uniqueIndex:unique_friendship;constraint:OnDelete:CASCADE"`
    User2ID   string    `json:"user2_id" gorm:"type:char(36);not null;index:idx_friendships_user2;uniqueIndex:unique_friendship;constraint:OnDelete:CASCADE"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (Friendship) TableName() string {
    return "friendships"
}