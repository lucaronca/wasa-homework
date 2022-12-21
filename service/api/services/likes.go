package services

import (
	"errors"

	"github.com/lucaronca/wasa-homework/service/api/models"
	"github.com/lucaronca/wasa-homework/service/api/repositories"
	"github.com/lucaronca/wasa-homework/service/globaltime"
)

var ErrNoLike = errors.New("Like not found")

// LikesService defines the api actions to like/unlike a photo
type LikesService interface {
	LikePhoto(int, int) error
	UnlikePhoto(int, int) error
	GetPhotoLikes(int, int) (*[]models.Like, error)
}

// likesService is a service that implements the logic for the LikesService
type likesService struct {
	ur repositories.UsersRepository
	br repositories.BansRepository
	lr repositories.LikesRepository
	pr repositories.PhotosRepository
}

// NewLikesService creates a default api service
func NewLikesService(ur repositories.UsersRepository, br repositories.BansRepository, lr repositories.LikesRepository, pr repositories.PhotosRepository) LikesService {
	return &likesService{
		ur: ur,
		br: br,
		lr: lr,
		pr: pr,
	}
}

// LikePhoto - Put a like to a photo
func (s *likesService) LikePhoto(photoId, userId int) error {
	photo, err := s.pr.GetPhotoById(photoId)
	if err != nil {
		return err
	}
	if photo == nil {
		return ErrNoPhoto
	}
	user, err := s.ur.GetUserById(userId)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNoUser
	}
	isBannedForUser, err := s.br.GetBanExists(userId, photo.Owner.Id)
	if err != nil {
		return err
	}
	if isBannedForUser {
		return ErrNoPhoto
	}
	isBannedForUser, err = s.br.GetBanExists(photo.Owner.Id, userId)
	if err != nil {
		return err
	}
	if isBannedForUser {
		return ErrNoPhoto
	}

	if err := s.lr.SetLike(photo.Id, user.Id, globaltime.Now()); err != nil {
		return err
	}
	return nil
}

// UnlikePhoto - Remove a like from a photo
func (s *likesService) UnlikePhoto(photoId, userId int) error {
	photo, err := s.pr.GetPhotoById(photoId)
	if err != nil {
		return err
	}
	if photo == nil {
		return ErrNoPhoto
	}
	user, err := s.ur.GetUserById(userId)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNoPhoto
	}
	isBannedForUser, err := s.br.GetBanExists(userId, photo.Owner.Id)
	if err != nil {
		return err
	}
	if isBannedForUser {
		return ErrNoPhoto
	}
	isBannedForUser, err = s.br.GetBanExists(photo.Owner.Id, userId)
	if err != nil {
		return err
	}
	if isBannedForUser {
		return ErrNoPhoto
	}

	likes, err := s.lr.GetLikes(s.ur.WithUsers(), s.ur.FilterByUserId(userId), s.pr.FilterByPhotoId(photoId))
	if err != nil {
		return err
	}
	if len(*likes) == 0 {
		return ErrNoLike
	}

	if err := s.lr.RemoveLike(photoId, userId); err != nil {
		return err
	}
	return nil
}

// GetPhotoLikes - Get photo likes
func (s *likesService) GetPhotoLikes(photoId, userId int) (*[]models.Like, error) {
	photo, err := s.pr.GetPhotoById(photoId)
	if err != nil {
		return nil, err
	}
	if photo == nil {
		return nil, ErrNoPhoto
	}
	user, err := s.ur.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNoUser
	}
	isBannedForUser, err := s.br.GetBanExists(userId, photo.Owner.Id)
	if err != nil {
		return nil, err
	}
	if isBannedForUser {
		return nil, ErrNoPhoto
	}
	isBannedForUser, err = s.br.GetBanExists(photo.Owner.Id, userId)
	if err != nil {
		return nil, err
	}
	if isBannedForUser {
		return nil, ErrNoPhoto
	}

	likes, err := s.lr.GetLikes(s.ur.WithUsers(), s.pr.FilterByPhotoId(photoId))
	if err != nil {
		return nil, err
	}

	if likes == nil || len(*likes) == 0 {
		empty := make([]models.Like, 0)
		return &empty, nil
	}
	return likes, nil
}
