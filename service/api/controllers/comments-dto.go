package controllers

import (
	"errors"
)

// CommentPhoto - Comment a photo
type CommentPhoto struct {
	// User name
	Content string `json:"content,omitempty"`
}

var ErrCommentContentIsZero = errors.New("Comment content is empty")
var ErrCommentContentIsNotValid = errors.New("Comment is too long")

// assertCommentPhotoValid checks if the required fields are not zero-ed
func assertCommentPhotoValid(obj CommentPhoto) error {
	length := len(obj.Content)
	switch {
	case length == 0:
		return ErrCommentContentIsZero
	case length > 500:
		return ErrCommentContentIsNotValid

	default:
		return nil
	}
}
