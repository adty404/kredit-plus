package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint   `gorm:"primarykey"`
	FullName  string `gorm:"type:varchar(255);not null"`
	Email     string `gorm:"type:varchar(100);unique;not null"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"type:varchar(50);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// HashPassword mengenkripsi password plain text menggunakan bcrypt.
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword memverifikasi apakah password yang diberikan cocok dengan hash.
func (u *User) CheckPassword(providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(providedPassword))
}
