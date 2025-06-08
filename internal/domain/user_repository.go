package domain

import "gorm.io/gorm"

type UserRepository interface {
	WithTx(tx *gorm.DB) UserRepository
	Save(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
}
