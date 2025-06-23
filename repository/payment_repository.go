package repository

import (
	"api.mijkomp.com/models/entity"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Save(db *gorm.DB, user entity.Payment) (entity.Payment, error)
}
