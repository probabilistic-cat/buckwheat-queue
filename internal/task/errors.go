package task

import "errors"

var (
	ErrorInvalidTask         = errors.New("Invalid task")
	ErrorInvalidStatus       = errors.New("Invalid status")
	ErrorInvalidStatusChange = errors.New("Invalid status change")
)
