package entity

import "time"

type UserVerification struct {
	Id           uint      `gorm:"primaryKey;"`
	UserId       uint      `gorm:"foreignKey;"`
	Code         string    `gorm:"type:varchar(6)"`
	AlreadyUsed  bool      `gorm:"not null"`
	CreatedById  uint      `gorm:"foreignKey; not null;"`
	CreatedAt    time.Time `gorm:"type:timestamptz; not null;"`
	ModifiedById uint      `gorm:"foreignKey; not null;"`
	ModifiedAt   time.Time `gorm:"type:timestamptz; not null;"`
}
