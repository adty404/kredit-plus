package usecase

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/adty404/kredit-plus/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestRegisterUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	usecase := NewUserUsecase(mockRepo)
	input := RegisterUserInput{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "consumer",
	}

	// Tentukan ekspektasi mock
	mockRepo.On("FindByEmail", input.Email).Return(nil, gorm.ErrRecordNotFound).Once()
	mockRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(nil).Once()

	// Act
	user, err := usecase.RegisterUser(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, input.Email, user.Email)
	assert.Equal(t, input.FullName, user.FullName)
	assert.Empty(t, user.Password, "Password should be empty in the response") // Pastikan hash tidak dikembalikan
	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_EmailAlreadyExists(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	usecase := NewUserUsecase(mockRepo)
	input := RegisterUserInput{Email: "test@example.com"}
	existingUser := &domain.User{Email: input.Email}

	mockRepo.On("FindByEmail", input.Email).Return(existingUser, nil).Once()

	// Act
	user, err := usecase.RegisterUser(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, fmt.Errorf("email '%s' already registered", input.Email), err)
	mockRepo.AssertExpectations(t)
}

// --- Test untuk LoginUser ---

func TestLoginUser_Success(t *testing.T) {
	// Arrange
	// Set JWT secret untuk testing
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	mockRepo := new(MockUserRepository)
	usecase := NewUserUsecase(mockRepo)
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	input := LoginInput{
		Email:    "test@example.com",
		Password: password,
	}

	existingUser := &domain.User{
		ID:       1,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "consumer",
	}

	mockRepo.On("FindByEmail", input.Email).Return(existingUser, nil).Once()

	// Act
	output, err := usecase.LoginUser(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.Token)
	mockRepo.AssertExpectations(t)
}

func TestLoginUser_InvalidCredentials_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	usecase := NewUserUsecase(mockRepo)
	input := LoginInput{Email: "notfound@example.com", Password: "password123"}

	mockRepo.On("FindByEmail", input.Email).Return(nil, gorm.ErrRecordNotFound).Once()

	// Act
	output, err := usecase.LoginUser(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, errors.New("invalid email or password"), err)
	mockRepo.AssertExpectations(t)
}

func TestLoginUser_InvalidCredentials_WrongPassword(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	usecase := NewUserUsecase(mockRepo)
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	input := LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	existingUser := &domain.User{
		ID:       1,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	mockRepo.On("FindByEmail", input.Email).Return(existingUser, nil).Once()

	// Act
	output, err := usecase.LoginUser(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, errors.New("invalid email or password"), err)
	mockRepo.AssertExpectations(t)
}
