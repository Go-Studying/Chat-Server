package service

import (
	"chat-server/internal/models"
	"chat-server/internal/repository"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	r *repository.UserRepository
}

func NewAuthService(r *repository.UserRepository) *AuthService {
	return &AuthService{r: r}
}

func (s AuthService) SignUp(email string, username string, password string) (uuid.UUID, error) {
	user, _ := s.r.FindUserByEmail(email)
	if user != nil {
		return uuid.Nil, fmt.Errorf("email already used")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user = &models.User{
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
	}

	id, err := s.r.SaveUser(user)
	return id, err
}

func (s AuthService) Login(email string, password string) (uuid.UUID, error) {
	user, _ := s.r.FindUserByEmail(email)
	if user == nil {
		return uuid.Nil, fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return uuid.Nil, fmt.Errorf("password error")
	}
	return user.ID, nil
}
