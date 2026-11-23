package task

import "slices"

type Status string

const (
	StatusPending   Status = "pending"
	StatusScheduled Status = "scheduled"
	StatusRunning   Status = "running"
	StatusRetrying  Status = "retrying"
	StatusSucceeded Status = "succeeded"
	StatusFailed    Status = "failed"
	StatusDead      Status = "dead"
	StatusCancelled Status = "cancelled"
)

var availableChange = map[Status][]Status{
	StatusPending:   {StatusRunning, StatusCancelled},
	StatusScheduled: {StatusPending, StatusCancelled},
	StatusRunning:   {StatusSucceeded, StatusFailed},
	StatusRetrying:  {StatusRunning, StatusCancelled},
	StatusSucceeded: {},
	StatusFailed:    {StatusRetrying, StatusDead},
	StatusDead:      {StatusRetrying},
	StatusCancelled: {},
}

func CanChange(from, to Status) bool {
	return slices.Contains(availableChange[from], to)
}
