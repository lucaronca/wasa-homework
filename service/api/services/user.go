package services

import (
	"errors"

	"github.com/lucaronca/wasa-homework/service/api/models"
	"github.com/lucaronca/wasa-homework/service/api/repositories"
)

var ErrNoUser = errors.New("User not found")
var ErrUserForbidden = errors.New("You are not authorized to do this operations")

type UsersService interface {
	GetUser(int, int) (*models.FullUser, error)
	GetUsers(int, string) (*[]models.BaseUser, error)
	UpdateUsername(int, string) (*models.FullUser, error)
}

type usersService struct {
	ur repositories.UsersRepository
	br repositories.BansRepository
	fr repositories.FollowsRepository
	pr repositories.PhotosRepository
}

func NewUsersService(
	ur repositories.UsersRepository,
	br repositories.BansRepository,
	fr repositories.FollowsRepository,
	pr repositories.PhotosRepository,
) UsersService {
	return &usersService{
		ur: ur,
		br: br,
		fr: fr,
		pr: pr,
	}
}

func (s *usersService) GetUser(userId int, targetUserId int) (*models.FullUser, error) {
	user, err := s.ur.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNoUser
	}
	targetUser, err := s.ur.GetUserById(targetUserId)
	if err != nil {
		return nil, err
	}
	if targetUser == nil {
		return nil, ErrNoUser
	}

	requestedUser := models.FullUser{}
	if userId == targetUserId {
		requestedUser.BannedForUser = nil
	} else {
		isBannedForUser, err := s.br.GetBanExists(userId, targetUserId)
		if err != nil {
			return nil, err
		}
		bfu := new(bool)
		if isBannedForUser {
			*bfu = true
			return &models.FullUser{
				BaseUser:      *targetUser,
				BannedForUser: bfu,
			}, nil
		} else {
			*bfu = false
			requestedUser.BannedForUser = bfu
		}

		isBannedForUser, err = s.br.GetBanExists(targetUserId, userId)
		if err != nil {
			return nil, err
		}
		if isBannedForUser {
			return nil, ErrNoUser
		}
	}

	fullUser, err := s.ur.GetFullUser(
		s.fr.WithTotalFollowers(),
		s.fr.WithTotalFollowings(),
		s.pr.WithTotalPhotos(),
		s.ur.FilterByUserId(targetUserId),
	)
	if err != nil {
		return nil, err
	}
	requestedUser.BaseUser = fullUser.BaseUser
	requestedUser.TotalFollowers = fullUser.TotalFollowers
	requestedUser.TotalFollowings = fullUser.TotalFollowings
	requestedUser.TotalPhotos = fullUser.TotalPhotos
	return &requestedUser, nil
}

func (s *usersService) GetUsers(userId int, username string) (*[]models.BaseUser, error) {
	users, err := s.ur.GetUsers(
		s.ur.FilterByUsername(username, false),
		s.br.WithoutBanners(userId),
	)
	if err != nil {
		return nil, err
	}
	if users == nil || len(*users) == 0 {
		empty := make([]models.BaseUser, 0)
		users = &empty
	}
	return users, nil
}

func (s *usersService) UpdateUsername(userId int, username string) (*models.FullUser, error) {
	user, err := s.ur.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNoUser
	}

	user.Username = username

	err = s.ur.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	fullUser, err := s.ur.GetFullUser(
		s.fr.WithTotalFollowers(),
		s.fr.WithTotalFollowings(),
		s.pr.WithTotalPhotos(),
		s.ur.FilterByUserId(userId),
	)
	return fullUser, err
}
