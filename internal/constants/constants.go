package constants

import "errors"

var ErrEntityNotFound = errors.New("entity not found")
var ErrInvalidRequest = errors.New("not a valid request")
var ErrEntityAlreadyExists = errors.New("entity already exists")
var ErrRequiredField = errors.New("field is required")

const DefaultError = "Internal Server Error"
