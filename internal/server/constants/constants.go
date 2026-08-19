package constants

import "errors"

var (
	ErrEntityNotFound      = errors.New("entity not found")
	ErrEntityAlreadyExists = errors.New("entity already exists")
	ErrRequiredField       = errors.New("field is required")
	ErrInvalidInput        = errors.New("invalid input")
)

const DefaultError = "Internal Server Error"
