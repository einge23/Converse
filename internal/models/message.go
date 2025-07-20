package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Metadata map[string]interface{}

func (m Metadata) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *Metadata) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	
	return json.Unmarshal(bytes, m)
}

type Message struct {
	MessageID   string     `json:"message_id" gorm:"column:message_id;type:char(36);primaryKey"`
	RoomID      *string    `json:"room_id" gorm:"column:room_id;type:char(36);index:idx_messages_room_id_created_at,priority:1;constraint:OnDelete:CASCADE"`
	ThreadID    *string    `json:"thread_id" gorm:"column:thread_id;type:char(36);index:idx_messages_thread_id_created_at,priority:1;constraint:OnDelete:CASCADE"`
	SenderID    *string    `json:"sender_id" gorm:"column:sender_id;type:char(36);index:idx_messages_sender_id;constraint:OnDelete:SET NULL"`
	ContentType string     `json:"content_type" gorm:"column:content_type;type:enum('text','image_url','file_url','system_notification','call_started','call_ended');not null;default:'text';index:idx_messages_content_type"`
	Content     string     `json:"content" gorm:"column:content;type:text;not null"`
	Metadata    *Metadata  `json:"metadata" gorm:"column:metadata;type:json"`
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at;not null;autoCreateTime;type:datetime(3);index:idx_messages_room_id_created_at,priority:2;index:idx_messages_thread_id_created_at,priority:2;index:idx_messages_created_at"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime(3)"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(3)"`
}

func (Message) TableName() string {
	return "messages"
}

func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	if m.MessageID == "" {
		m.MessageID = uuid.New().String()
	}
	return nil
}

type MessageWithAttachments struct {
    Message
    Attachments []MessageAttachment `json:"attachments"`
}