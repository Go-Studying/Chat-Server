package service

import (
	"chat-server/internal/models"
	"chat-server/internal/repository"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicateKey = errors.New("email already used")

type AuthService struct {
	r *repository.UserRepository
}

func NewAuthService(r *repository.UserRepository) *AuthService {
	return &AuthService{r: r}
}

func (s AuthService) SignUp(email string, username string, password string) (uuid.UUID, error) {
	user, err := s.r.FindUserByEmail(email)
	if user != nil {
		return uuid.Nil, ErrDuplicateKey
	}
	if err != nil {
		return uuid.Nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}

	user = &models.User{
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
	}

	id, err := s.r.SaveUser(user)
	return id, err
}

func (s AuthService) Login(email string, password string) (uuid.UUID, error) {
	user, err := s.r.FindUserByEmail(email)
	if user == nil {
		return uuid.Nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return uuid.Nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return uuid.Nil, fmt.Errorf("password error")
	}
	return user.ID, nil
}
