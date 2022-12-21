package services

import (
	"errors"

	"github.com/lucaronca/wasa-homework/service/api/models"
	"github.com/lucaronca/wasa-homework/service/api/repositories"
	"github.com/lucaronca/wasa-homework/service/globaltime"
)

var ErrNoComment = errors.New("Comment not found")
var ErrDeleteNotAllowed = errors.New("You can't delete this comment")

// CommentsService defines the api actions to comment/uncomment a photo
type CommentsService interface {
	CommentPhoto(int, int, string) (*models.Comment, error)
	UncommentPhoto(int, int, int) error
	GetPhotoComments(int, int) (*[]models.Comment, error)
}

// commentsService is a service that implements the logic for the CommentsService
type commentsService struct {
	ur repositories.UsersRepository
	br repositories.BansRepository
	cr repositories.CommentsRepository
	pr repositories.PhotosRepository
}

// NewCommentsService creates a default api service
func NewCommentsService(ur repositories.UsersRepository, br repositories.BansRepository, cr repositories.CommentsRepository, pr repositories.PhotosRepository) CommentsService {
	return &commentsService{
		ur: ur,
		br: br,
		cr: cr,
		pr: pr,
	}
}

// CommentPhoto - Put a comment to a photo
func (s *commentsService) CommentPhoto(photoId, userId int, content string) (*models.Comment, error) {
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

	commentId, err := s.cr.SetComment(photo.Id, user.Id, globaltime.Now(), content)
	if err != nil {
		return nil, err
	}
	comment, err := s.cr.GetCommentById(commentId, s.ur.WithUsers())
	if err != nil {
		return nil, err
	}
	return comment, err
}

// UncommentPhoto - Remove a comment from a photo
func (s *commentsService) UncommentPhoto(photoId, commentId, userId int) error {
	photo, err := s.pr.GetPhotoById(photoId)
	if err != nil {
		return err
	}
	if photo == nil {
		return ErrNoPhoto
	}
	comment, err := s.cr.GetCommentById(commentId, s.ur.WithUsers())
	if err != nil {
		return err
	}
	if comment == nil {
		return ErrNoComment
	}
	if comment.Owner.Id != userId {
		return ErrDeleteNotAllowed
	}
	user, err := s.ur.GetUserById(userId)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNoUser
	}

	if err := s.cr.RemoveComment(commentId); err != nil {
		return err
	}
	return nil
}

// GetPhotoComments - Get photo comments
func (s *commentsService) GetPhotoComments(photoId, userId int) (*[]models.Comment, error) {
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

	comments, err := s.cr.GetComments(s.ur.WithUsers(), s.pr.FilterByPhotoId(photoId))
	if err != nil {
		return nil, err
	}

	if comments == nil || len(*comments) == 0 {
		empty := make([]models.Comment, 0)
		return &empty, nil
	}
	return comments, nil
}
