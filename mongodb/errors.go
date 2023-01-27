package mongodb

import "errors"

var (
	TooManyResultsError   = errors.New("too many results detected, requires investigation")
	DocumentNotFoundError = errors.New("model not found")
	DocumentExistsError   = errors.New("model already present")
)
