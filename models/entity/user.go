package entity

import "time"

type User struct {
	Id          uint      `gorm:"primaryKey;autoincrement"`
	UserName    string    `gorm:"type:varchar(128); not null"`
	FullName    string    `gorm:"type:varchar(128); not null"`
	Email       string    `gorm:"type:varchar(128); unique; not null"`
	Password    *string   `gorm:"type:varchar(128); null"`
	CreatedById uint      `gorm:"type:bigint; not null"`
	CreatedAt   time.Time `gorm:"type:timestamptz; not null"`
}
