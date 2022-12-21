package models

import (
	"time"
)

// Photo - A photo published by a user
type Photo struct {

	// Unique identifier of a Photo
	Id int `json:"id,omitempty"`

	// Image URL
	Url string `json:"url,omitempty"`

	// Image likes number
	TotalLikes int `json:"totalLikes"`

	// Image comments number
	TotalComments int `json:"totalComments"`

	// Image upload date
	UploadDate time.Time `json:"uploadDate,omitempty"`

	Owner BaseUser `json:"owner,omitempty"`
}
