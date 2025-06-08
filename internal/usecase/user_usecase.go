package usecase

import (
	"errors"
	"fmt"
	"github.com/adty404/kredit-plus/internal/auth"
	"github.com/adty404/kredit-plus/internal/domain"
	"gorm.io/gorm"
)

type UserUsecase interface {
	RegisterUser(input RegisterUserInput) (*domain.User, error)
	LoginUser(input LoginInput) (*LoginOutput, error)
}

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (uc *userUsecase) RegisterUser(input RegisterUserInput) (*domain.User, error) {
	// Cek apakah email sudah ada
	_, err := uc.userRepo.FindByEmail(input.Email)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("email '%s' already registered", input.Email)
	}

	// Buat user baru
	newUser := &domain.User{
		FullName: input.FullName,
		Email:    input.Email,
		Role:     input.Role,
	}

	// Hash password sebelum disimpan
	if err := newUser.HashPassword(input.Password); err != nil {
		return nil, err
	}

	// Simpan user ke database
	if err := uc.userRepo.Save(newUser); err != nil {
		return nil, err
	}

	// Jangan kembalikan hash password di respons
	newUser.Password = ""

	return newUser, nil
}

func (uc *userUsecase) LoginUser(input LoginInput) (*LoginOutput, error) {
	// Cari user berdasarkan email
	user, err := uc.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Verifikasi password
	if err := user.CheckPassword(input.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Jika password cocok, buat token JWT
	token, err := auth.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{Token: token}, nil
}
