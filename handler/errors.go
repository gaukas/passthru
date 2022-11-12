package handler

import "errors"

var (
	ErrServerStopped = errors.New("server stopped")
	ErrUnknownAction = errors.New("unknown action")
)
