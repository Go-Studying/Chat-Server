package service

import "errors"

// Auth
var ErrDuplicateKey = errors.New("email already used")

// ChatRoom
var ErrRoomNotFound = errors.New("room not found")
