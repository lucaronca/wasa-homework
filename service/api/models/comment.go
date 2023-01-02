package models

import "time"

// Comment - A comment to a photo published by a user
type Comment struct {

	// Unique identifier of a Photo comment
	Id int `json:"id,omitempty"`

	// Like date
	Date time.Time `json:"date,omitempty"`

	// Comment's content
	Content string `json:"content,omitempty"`

	Photo BasePhoto `json:"photo,omitempty"`

	Owner BaseUser `json:"owner,omitempty"`
}
