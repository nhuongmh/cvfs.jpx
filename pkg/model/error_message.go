package model

import "errors"

var (
	ErrNotImplemented          = errors.New("method not implemented")
	ErrNoData                  = errors.New("no data")
	ErrServiceIsNotInitialized = errors.New("service is not initialized")
)
