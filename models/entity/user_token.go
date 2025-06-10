package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserToken struct {
	Token     uuid.UUID `gorm:"type:uuid; primaryKey"`
	UserId    uint      `gorm:"foreignKey;required"`
	IpAddress string    `gorm:"type:varchar(200);required"`
	CreatedAt time.Time `gorm:"type:timestamptz; not null;"`
}
