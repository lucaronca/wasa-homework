package services

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
	"github.com/lucaronca/wasa-homework/service/api/models"
	"github.com/lucaronca/wasa-homework/service/api/repositories"
	"github.com/lucaronca/wasa-homework/service/globaltime"
)

var ErrNoPhoto = errors.New("Photo not found")

// PhotosService defines the api actions to manage photos
type PhotosService interface {
	GetUserPhotos(int, int, int, int) (*models.PaginatedPhotos, error)
	GetStream(int, int, int) (*models.PaginatedPhotos, error)
	CreatePhoto(int, []byte, string) (*models.Photo, error)
	DeletePhoto(int, int) error
}

// photosService is a service that implements the logic for the PhotosService
type photosService struct {
	photosDirectory string
	photosUrlPath   string
	ur              repositories.UsersRepository
	br              repositories.BansRepository
	pr              repositories.PhotosRepository
	lr              repositories.LikesRepository
	cr              repositories.CommentsRepository
	fr              repositories.FollowsRepository
}

// NewPhotosService creates a default api service
func NewPhotosService(
	photosDirectory string,
	photosUrlPath string,
	ur repositories.UsersRepository,
	br repositories.BansRepository,
	pr repositories.PhotosRepository,
	lr repositories.LikesRepository,
	cr repositories.CommentsRepository,
	fr repositories.FollowsRepository,
) PhotosService {
	return &photosService{
		photosDirectory: photosDirectory,
		photosUrlPath:   photosUrlPath,
		ur:              ur,
		br:              br,
		pr:              pr,
		lr:              lr,
		cr:              cr,
		fr:              fr,
	}
}

func (s *photosService) GetUserPhotos(userId, targetUserId, offset, limit int) (*models.PaginatedPhotos, error) {
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
	if userId != targetUserId {
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
	}

	out := NewWorkersFacade(
		NewJob(func(sendRes SendFunc) {
			result, err := s.pr.GetPhotos(
				offset,
				limit,
				s.ur.WithUsers(),
				s.lr.WithTotalLikes(),
				s.cr.WithTotalComments(),
				s.ur.FilterByUserId(targetUserId),
			)
			if err != nil {
				sendRes(nil, err)
				return
			}
			if result == nil || len(*result) == 0 {
				empty := make([]models.Photo, 0)
				result = &empty
			}
			sendRes(result, nil)
		}),
		NewJob(func(sendRes SendFunc) {
			result, err := s.pr.GetPhotosCount(s.ur.FilterByUserId(targetUserId))
			if err != nil {
				sendRes(nil, err)
			}
			sendRes(result, nil)
		}),
	)

	var entries *[]models.Photo
	var totalCount int
	for work := range out {
		if work.err != nil {
			return nil, work.err
		}

		switch work.idx {
		case 0:
			entries, _ = work.res.(*[]models.Photo)
		case 1:
			totalCount, _ = work.res.(int)
		}
	}

	return &models.PaginatedPhotos{
		Offset:     offset,
		Limit:      limit,
		Entries:    entries,
		TotalCount: totalCount,
	}, nil
}

func (s *photosService) GetStream(targetUserId, offset, limit int) (*models.PaginatedPhotos, error) {
	user, err := s.ur.GetUserById(targetUserId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNoUser
	}

	out := NewWorkersFacade(
		NewJob(func(sendRes SendFunc) {
			result, err := s.pr.GetPhotos(
				offset,
				limit,
				s.ur.WithUsers(),
				s.lr.WithTotalLikes(),
				s.cr.WithTotalComments(),
				s.fr.FilterByFollowerId(targetUserId),
				s.br.WithoutBanned(targetUserId),
			)
			if err != nil {
				sendRes(nil, err)
				return
			}
			if result == nil || len(*result) == 0 {
				empty := make([]models.Photo, 0)
				result = &empty
			}
			sendRes(result, nil)
		}),
		NewJob(func(sendRes SendFunc) {
			result, err := s.pr.GetPhotosCount(
				s.fr.FilterByFollowerId(targetUserId),
				s.br.WithoutBanned(targetUserId),
			)
			if err != nil {
				sendRes(nil, err)
			}
			sendRes(result, nil)
		}),
	)

	var entries *[]models.Photo
	var totalCount int
	for work := range out {
		if work.err != nil {
			return nil, work.err
		}

		switch work.idx {
		case 0:
			entries, _ = work.res.(*[]models.Photo)
		case 1:
			totalCount, _ = work.res.(int)
		}
	}

	return &models.PaginatedPhotos{
		Offset:     offset,
		Limit:      limit,
		Entries:    entries,
		TotalCount: totalCount,
	}, nil
}

func (s *photosService) CreatePhoto(userId int, photo []byte, ext string) (*models.Photo, error) {
	photoName, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	photoNameWithExt := photoName.String() + "." + ext
	photoFilePath := filepath.Join(s.photosDirectory, photoNameWithExt)

	photoId, err := s.pr.SetPhoto(filepath.Join(s.photosUrlPath, photoNameWithExt), userId, globaltime.Now())
	if err != nil {
		return nil, err
	}

	file, err := os.Create(photoFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := bytes.NewReader(photo)
	_, err = io.Copy(file, r)
	if err != nil {
		return nil, err
	}

	newPhoto, err := s.pr.GetPhotoById(photoId)
	if err != nil {
		return nil, err
	}
	return newPhoto, nil
}

// DeletePhoto - Delete a photos
func (s *photosService) DeletePhoto(userId, photoId int) error {
	user, err := s.ur.GetUserById(userId)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNoUser
	}
	photo, err := s.pr.GetPhotoById(photoId)
	if err != nil {
		return err
	}
	if photo == nil {
		return ErrNoPhoto
	}
	if userId != photo.Owner.Id {
		return ErrUserForbidden
	}

	if err := s.pr.RemovePhoto(photoId); err != nil {
		return err
	}
	filePath := filepath.Join(s.photosDirectory, filepath.Base(photo.Url))
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}
