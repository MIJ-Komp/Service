package repository_impl

import (
	"api.mijkomp.com/models/entity"
	"gorm.io/gorm"
)

type PaymentRepositoryImpl struct{}

func NewPaymentRepository() *PaymentRepositoryImpl {
	return &PaymentRepositoryImpl{}
}

func (repository *PaymentRepositoryImpl) Save(db *gorm.DB, payment entity.Payment) (entity.Payment, error) {
	err := db.Save(&payment).Error

	return payment, err
}
