package models

type MessageAttachment struct {
    ID        string `gorm:"column:id;primaryKey;type:varchar(36)" json:"id"`
    MessageID string `gorm:"column:message_id;type:varchar(36);not null;index" json:"message_id"`
    FileURL   string `gorm:"column:file_url;type:varchar(500);not null" json:"file_url"`
    FileKey   string `gorm:"column:file_key;type:varchar(200);not null" json:"file_key"`
    FileType  string `gorm:"column:file_type;type:varchar(50);not null" json:"file_type"`
    FileSize  *int   `gorm:"column:file_size;type:int" json:"file_size"`
    Filename  *string `gorm:"column:filename;type:varchar(255)" json:"filename"`
}

func (MessageAttachment) TableName() string {
	return "message_attachments"
}