package main

import "errors"

var ErrCategoryNotFound = errors.New("category not found")
var ErrPostNotFound = errors.New("post not found")
var ErrUserNotLoggedIn = errors.New("user not logged in")
