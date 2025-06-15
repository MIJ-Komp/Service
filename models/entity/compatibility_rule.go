package entity

import "time"

type CompatibilityRule struct {
	Id                      uint      `gorm:"primaryKey; type:bigint; autoincrement;"`
	SourceComponentTypeCode string    `gorm:"type:varchar(256); not null"`
	TargetComponentTypeCode string    `gorm:"type:varchar(256); not null"`
	SourceKey               string    `gorm:"type:varchar(256); not null;"`
	TargetKey               string    `gorm:"type:varchar(256); not null;"`
	Condition               string    `gorm:"type:varchar(10); not null;"`
	ValueType               string    `gorm:"type:varchar(24); not null;"`
	ErrorMessage            string    `gorm:"type:varchar(256); not null;"`
	IsActive                bool      `gorm:"type:bool; not null;"`
	CreatedById             uint      `gorm:"type:bigint; not null"`
	CreatedAt               time.Time `gorm:"type:timestamptz; not null"`
	ModifiedById            uint      `gorm:"type:bigint; not null"`
	ModifiedAt              time.Time `gorm:"type:timestamptz; not null"`

	CreatedBy  User
	ModifiedBy User
}
