package controllers

import (
	"errors"
	"regexp"
)

// DoLoginRequest - Name of the user
type DoLoginRequest struct {

	// User name
	Name string `json:"name,omitempty"`
}

var ErrLoginNameIsZero = errors.New("Login name is zero value")
var ErrLoginNameIsNotValid = errors.New("Name should be at least 3 characters long")

func parseUsernameParameter(param string) error {
	match, err := regexp.MatchString(`^.*?$`, param)

	if err != nil || !match {
		return ErrLoginNameIsNotValid
	}

	return nil
}

// assertDoLoginRequestValid checks if the required fields are not zero-ed
func assertDoLoginRequestValid(obj DoLoginRequest) error {
	switch len(obj.Name) {
	case 0:
		return ErrLoginNameIsZero
	case 1, 2:
		return ErrLoginNameIsNotValid
	}

	return parseUsernameParameter(obj.Name)
}

// DoLoginResponse - Token identifier for a logged in user
type DoLoginResponse struct {

	// Token identifier for a logged in user
	Identifier string `json:"identifier,omitempty"`
}
