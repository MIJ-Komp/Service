package entity

import "time"

type Menu struct {
	Id           uint      `gorm:"primaryKey;autoincrement"`
	Name         string    `gorm:"type:varchar(128); not null"`
	ParentId     *uint     `gorm:"type:bigint;null"`
	Path         *string   `gorm:"type:varchar(256); null"`
	CreatedById  uint      `gorm:"type:bigint; not null"`
	CreatedAt    time.Time `gorm:"type:timestamptz; not null"`
	ModifiedById uint      `gorm:"type:bigint; not null"`
	ModifiedAt   time.Time `gorm:"type:timestamptz; not null"`

	CreatedBy  User
	ModifiedBy User

	MenuItems []MenuItem `gorm:"foreignKey:menu_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type MenuItem struct {
	Id                uint `gorm:"primaryKey;autoincrement"`
	MenuId            uint `gorm:"type:bigint"`
	ProductCategoryId uint `gorm:"type:bigint"`

	ProductCategory ProductCategory `gorm:"foreignKey:product_category_id;references:Id;"`
}
