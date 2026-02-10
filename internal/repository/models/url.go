package models

import (
	"github.com/google/uuid"
)

type Url struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UrlCode     string    `gorm:"type:VARCHAR(10)"`
	OriginalUrl string    `gorm:"type:VARCHAR(2048)"`
}

func (Url) TableName() string {
	return "urls"
}
