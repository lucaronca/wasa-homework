package controllers

import (
	"errors"
)

// CommentPhoto - Comment a photo
type CommentPhoto struct {
	// User name
	Content string `json:"content,omitempty"`
}

var ErrCommentContentIsZero = errors.New("Comment content is zero value")
var ErrCommentConentIsNotValid = errors.New("Comment is too long")

func parseContentParameter(param string) error {
	if len(param) > 1000 {
		return ErrCommentConentIsNotValid
	}

	return nil
}

// assertCommentPhotoValid checks if the required fields are not zero-ed
func assertCommentPhotoValid(obj CommentPhoto) error {
	if len(obj.Content) == 0 {
		return ErrCommentContentIsZero
	}

	return parseContentParameter(obj.Content)
}
