package task

import (
	"encoding/json"
	"time"
	"uuid"
)

type Task struct {
	ID          uuid.UUID
	Type        string
	Payload     json.RawMessage
	Status      Status
	Attempts    int
	MaxAttempts int
	Timeout     time.Duration
	ScheduledAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
	Result      json.RawMessage
	Error       *string
}

func (t Task) Validate() error {
	if t.Type == "" ||
		!json.Valid(t.Payload) ||
		t.MaxAttempts < 1 ||
		t.Attempts > t.MaxAttempts {
		return ErrorInvalidTask
	}

	return nil
}
