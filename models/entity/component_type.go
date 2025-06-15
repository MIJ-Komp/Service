package entity

import "time"

type ComponentType struct {
	Id           uint      `gorm:"primaryKey;autoincrement"`
	Code         string    `gorm:"type:varchar(256); not null"`
	Name         string    `gorm:"type:varchar(1024); not null"`
	Description  string    `gorm:"type:varchar(4096); not null"`
	CreatedById  uint      `gorm:"type:bigint; not null"`
	CreatedAt    time.Time `gorm:"type:timestamptz; not null"`
	ModifiedById uint      `gorm:"type:bigint; not null"`
	ModifiedAt   time.Time `gorm:"type:timestamptz; not null"`

	CreatedBy  User
	ModifiedBy User
}
