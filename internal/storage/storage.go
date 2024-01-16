package storage

import "errors"

var (
	ErrorUserNotFound = errors.New("url not found")
	ErrorURLExists    = errors.New("url exists")
)
