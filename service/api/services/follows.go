package services

import (
	"github.com/lucaronca/wasa-homework/service/api/models"
	"github.com/lucaronca/wasa-homework/service/api/repositories"
)

// FollowsService defines the api actions to follow/unfollow a user
type FollowsService interface {
	FollowUser(int, int) error
	GetUserFollowers(int, int) (*[]models.BaseUser, error)
	GetUserFollowings(int, int) (*[]models.BaseUser, error)
	UnfollowUser(int, int) error
}

// followsService is a service that implements the logic for the FollowsService
type followsService struct {
	ur repositories.UsersRepository
	br repositories.BansRepository
	fr repositories.FollowsRepository
}

// NewFollowsService creates a default api service
func NewFollowsService(ur repositories.UsersRepository, br repositories.BansRepository, fr repositories.FollowsRepository) FollowsService {
	return &followsService{
		ur: ur,
		br: br,
		fr: fr,
	}
}

// FollowUser - Follow a user
func (s *followsService) FollowUser(followerUserId, followingUserId int) error {
	user, err := s.ur.GetUserById(followerUserId)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNoUser
	}
	user, err = s.ur.GetUserById(followingUserId)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNoUser
	}
	isBannedForUser, err := s.br.GetBanExists(followerUserId, followingUserId)
	if err != nil {
		return err
	}
	if isBannedForUser {
		return ErrNoUser
	}
	isBannedForUser, err = s.br.GetBanExists(followingUserId, followerUserId)
	if err != nil {
		return err
	}
	if isBannedForUser {
		return ErrNoUser
	}

	if err := s.fr.SetFollow(followerUserId, followingUserId); err != nil {
		return err
	}
	return nil
}

// GetUserFollowers - Get user followers
func (s *followsService) GetUserFollowers(userId int, targetUserId int) (*[]models.BaseUser, error) {
	user, err := s.ur.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNoUser
	}
	user, err = s.ur.GetUserById(targetUserId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNoUser
	}
	isBannedForUser, err := s.br.GetBanExists(userId, targetUserId)
	if err != nil {
		return nil, err
	}
	if isBannedForUser {
		return nil, ErrNoUser
	}
	isBannedForUser, err = s.br.GetBanExists(targetUserId, userId)
	if err != nil {
		return nil, err
	}
	if isBannedForUser {
		return nil, ErrNoUser
	}

	users, err := s.ur.GetUsers(s.fr.FilterByFollowingId(targetUserId))
	if err != nil {
		return nil, err
	}

	if users == nil || len(*users) == 0 {
		empty := make([]models.BaseUser, 0)
		return &empty, nil
	}
	return users, nil
}

// GetUserFollowings - Get user followings
func (s *followsService) GetUserFollowings(userId int, targetUserId int) (*[]models.BaseUser, error) {
	user, err := s.ur.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNoUser
	}
	user, err = s.ur.GetUserById(targetUserId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNoUser
	}
	isBannedForUser, err := s.br.GetBanExists(userId, targetUserId)
	if err != nil {
		return nil, err
	}
	if isBannedForUser {
		return nil, ErrNoUser
	}
	isBannedForUser, err = s.br.GetBanExists(targetUserId, userId)
	if err != nil {
		return nil, err
	}
	if isBannedForUser {
		return nil, ErrNoUser
	}

	users, err := s.ur.GetUsers(s.fr.FilterByFollowerId(targetUserId))
	if err != nil {
		return nil, err
	}
	if users == nil || len(*users) == 0 {
		empty := make([]models.BaseUser, 0)
		return &empty, nil
	}
	return users, nil
}

// UnfollowUser - Unfollow a user
func (s *followsService) UnfollowUser(followerUserId, followingUserId int) error {
	user, err := s.ur.GetUserById(followerUserId)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNoUser
	}
	user, err = s.ur.GetUserById(followingUserId)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNoUser
	}
	isBannedForUser, err := s.br.GetBanExists(followerUserId, followingUserId)
	if err != nil {
		return err
	}
	if isBannedForUser {
		return ErrNoUser
	}
	isBannedForUser, err = s.br.GetBanExists(followingUserId, followerUserId)
	if err != nil {
		return err
	}
	if isBannedForUser {
		return ErrNoUser
	}

	if err := s.fr.RemoveFollow(followerUserId, followingUserId); err != nil {
		return err
	}
	return nil
}
