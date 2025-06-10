package entity

import "time"

type ProductCategory struct {
	Id           uint      `gorm:"primaryKey;autoincrement"`
	Name         string    `gorm:"type:varchar(1024); not null"`
	ParentId     *uint     `gorm:"type:bigint;null"`
	CreatedById  uint      `gorm:"type:bigint; not null"`
	CreatedAt    time.Time `gorm:"type:timestamptz; not null"`
	ModifiedById uint      `gorm:"type:bigint; not null"`
	ModifiedAt   time.Time `gorm:"type:timestamptz; not null"`

	CreatedBy  User
	ModifiedBy User
}
