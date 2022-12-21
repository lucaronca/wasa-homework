package services

import (
	"crypto/sha1"
	"encoding/hex"
	"time"

	"github.com/lucaronca/wasa-homework/service/api/models"
	"github.com/lucaronca/wasa-homework/service/api/repositories"
)

type AuthService interface {
	DoLogin(username string) (token string, newUser bool, err error)
	Authorize(token string) (*models.BaseUser, error)
}

type authService struct {
	ar repositories.AuthRepository
	ur repositories.UsersRepository
}

func NewAuthService(ar repositories.AuthRepository, ur repositories.UsersRepository) AuthService {
	return &authService{
		ar: ar,
		ur: ur,
	}
}

func (a *authService) generateToken(username string) string {
	now := time.Now().String()
	hasher := sha1.New()
	hasher.Write([]byte(username + now))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (s *authService) DoLogin(username string) (string, bool, error) {
	user, err := s.ur.GetUser(s.ur.FilterByUsername(username, true))
	if err != nil {
		return "", false, err
	}
	// User doesn't exist, create it and its token
	if user == nil {
		userId, err := s.ur.CreateUser(&models.BaseUser{Username: username})
		if err != nil {
			return "", false, err
		}
		token := s.generateToken(username)
		if err := s.ar.SetToken(userId, token); err != nil {
			return "", false, err
		}
		return token, true, nil
	}

	// User exists, just return its token
	token, err := s.ar.GetToken(s.ur.FilterByUserId(user.Id))
	if err != nil {
		return "", false, err
	}
	return token, false, nil
}

func (s *authService) Authorize(token string) (*models.BaseUser, error) {
	user, err := s.ur.GetUser(s.ar.WithTokens(), s.ar.FilterByToken(token))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNoUser
	}
	return user, nil
}
