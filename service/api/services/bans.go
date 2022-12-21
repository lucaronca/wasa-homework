package services

import (
	"github.com/lucaronca/wasa-homework/service/api/repositories"
)

// BansService defines the api actions to ban/unban a user
type BansService interface {
	BanUser(int, int) error
	UnbanUser(int, int) error
	IsBannedForUser(int, int) (bool, error)
}

// bansService is a service that implements the logic for the BansService
type bansService struct {
	ur repositories.UsersRepository
	br repositories.BansRepository
	fr repositories.FollowsRepository
}

// NewBansService creates a default api service
func NewBansService(
	ur repositories.UsersRepository,
	br repositories.BansRepository,
	fr repositories.FollowsRepository,
) BansService {
	return &bansService{
		ur: ur,
		br: br,
		fr: fr,
	}
}

// BanUser - Ban a user
func (s *bansService) BanUser(userId, bannedId int) error {
	user, err := s.ur.GetUserById(userId)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNoUser
	}
	user, err = s.ur.GetUserById(bannedId)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNoUser
	}
	if err := s.br.SetBan(userId, bannedId); err != nil {
		return err
	}
	if err = s.fr.RemoveFollow(userId, bannedId); err != nil {
		return err
	}
	if err = s.fr.RemoveFollow(bannedId, userId); err != nil {
		return err
	}
	return nil
}

// UnbanUser - Unban a user
func (s *bansService) UnbanUser(userId, bannedId int) error {
	user, err := s.ur.GetUserById(userId)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNoUser
	}
	user, err = s.ur.GetUserById(bannedId)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNoUser
	}

	if err := s.br.RemoveBan(userId, bannedId); err != nil {
		return err
	}
	return nil
}

// IsBannedForUser - If a target user is banned/unbanned a user
func (s *bansService) IsBannedForUser(userId, targetUserId int) (bool, error) {
	user, err := s.ur.GetUserById(userId)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, ErrNoUser
	}
	user, err = s.ur.GetUserById(targetUserId)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, ErrNoUser
	}

	banExists, err := s.br.GetBanExists(userId, targetUserId)
	if err != nil {
		return false, err
	}
	return banExists, nil
}
