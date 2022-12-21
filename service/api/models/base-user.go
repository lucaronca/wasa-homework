package models

// BaseUser - Provides basic information about someone with a WASA Photo account.
type BaseUser struct {

	// Unique identifier of a user
	Id int `json:"id"`

	// User name
	Username string `json:"username"`
}
