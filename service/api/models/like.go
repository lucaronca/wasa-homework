package models

import "time"

// Like - A like to a photo published by a user
type Like struct {
	Id int `json:"id"`

	// Like date
	Date time.Time `json:"uploadDate,omitempty"`

	Photo BasePhoto `json:"photo,omitempty"`

	Owner BaseUser `json:"owner,omitempty"`
}
