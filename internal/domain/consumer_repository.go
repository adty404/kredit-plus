package domain

import "gorm.io/gorm"

type (
	ConsumerRepository interface {
		WithTx(tx *gorm.DB) ConsumerRepository
		FindByIDForUpdate(id uint) (*Consumer, error)
		Save(consumer *Consumer) error
		Update(id uint, updates map[string]interface{}) error
		FindByUserID(userID uint) (*Consumer, error)
		FindByID(id uint) (*Consumer, error)
		FindByNIK(nik string) (*Consumer, error)
		FindAll() ([]*Consumer, error)
		Delete(id uint) error
	}
)
