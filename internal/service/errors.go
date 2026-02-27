package service

import "errors"

// Auth
var ErrDuplicateKey = errors.New("email already used")

// ChatRoom
var ErrRoomNotFound = errors.New("room not found")
var ErrNotOwner = errors.New("only owner can delete room")
var ErrNotMember = errors.New("not a member of this room")
var ErrMemberExists = errors.New("member already exists")
