package controllers

import (
	"errors"
)

var ErrSeMyUserNameOpIsZero = errors.New("Op is zero value")
var ErrSeMyUserNameOpIsNotValid = errors.New("Op value is not valid")
var ErrSeMyUserNamePathIsZero = errors.New("Path is zero value")
var ErrSeMyUserNamePathIsNotValid = errors.New("Path value is not valid")
var ErrSetMyUserNameValueIsZero = errors.New("Value name is zero value")
var ErrSetMyUserNameValueIsNotValid = errors.New("Value should be at least 3 characters long")

type SetMyUserNameRequest struct {
	// Patch operation type
	Op string `json:"op,omitempty"`

	// Field to patch
	Path string `json:"path,omitempty"`

	// New field valued
	Value string `json:"value,omitempty"`
}

// assertSetMyUserNameRequestValid checks if the required fields are not zero-ed
func assertSetMyUserNameRequestValid(obj SetMyUserNameRequest) error {
	if obj.Op == "" {
		return ErrSeMyUserNameOpIsZero
	}
	if obj.Op != "replace" {
		return ErrSeMyUserNameOpIsNotValid
	}
	if obj.Path == "" {
		return ErrSeMyUserNamePathIsZero
	}
	if obj.Path != "/username" {
		return ErrSeMyUserNamePathIsNotValid
	}

	switch len(obj.Value) {
	case 0:
		return ErrSetMyUserNameValueIsZero
	case 1, 2:
		return ErrSetMyUserNameValueIsNotValid
	}

	return parseUsernameParameter(obj.Value)
}
