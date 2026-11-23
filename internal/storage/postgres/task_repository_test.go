package postgres_test

import (
	"buckwheat-queue/internal/storage/postgres"
	"buckwheat-queue/internal/task"
	"buckwheat-queue/internal/test"
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"
	"uuid"
)

func TestTaskRepository(t *testing.T) {
	ctx := t.Context()
	cfg := test.GetConfig(t)
	db := test.GetDB(t, cfg)
	repo := postgres.NewTaskRepository(db)

	now := time.Now().Truncate(time.Microsecond)
	id := uuid.NewV7()

	t.Cleanup(func() {
		_, err := db.Pool.Exec(context.Background(), `DELETE FROM task WHERE id = $1`, id)

		if err != nil {
			t.Errorf("Cleanup task: %v", err)
		}
	})

	original := &task.Task{
		ID:          id,
		Type:        "email.send",
		Payload:     json.RawMessage(`{"to":"user@example.com"}`),
		Status:      task.StatusPending,
		Attempts:    0,
		MaxAttempts: 3,
		Timeout:     5 * time.Minute,
		ScheduledAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Create
	err := repo.Create(ctx, original)
	if err != nil {
		t.Fatalf("Create task: %v", err)
	}

	// Exists
	exists, err := repo.Exists(ctx, original.ID)
	if err != nil {
		t.Fatalf("Check task exists: %v", err)
	}
	if !exists {
		t.Errorf("Task not exists")
	}

	// Read
	actual, err := repo.GetByID(ctx, original.ID)
	if err != nil {
		t.Fatalf("Get created task: %v", err)
	}

	// Compare fields.
	checkFields(t, actual.ID, original.ID, "ID")
	checkFields(t, actual.Type, original.Type, "Type")
	checkJSONFields(t, actual.Payload, original.Payload, "Payload")
	checkFields(t, actual.Status, original.Status, "Status")
	checkFields(t, actual.Attempts, original.Attempts, "Attempts")
	checkFields(t, actual.MaxAttempts, original.MaxAttempts, "MaxAttempts")
	checkFields(t, actual.Timeout, original.Timeout, "Timeout")
	checkFields(t, actual.ScheduledAt, original.ScheduledAt, "ScheduledAt")
	checkFields(t, actual.CreatedAt, original.CreatedAt, "CreatedAt")
	checkFields(t, actual.UpdatedAt, original.UpdatedAt, "UpdatedAt")

	// Delete
	err = repo.Delete(ctx, actual.ID)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	exists, err = repo.Exists(ctx, actual.ID)
	if err != nil {
		t.Fatalf("Check task exists: %v", err)
	}
	if exists {
		t.Errorf("Task exists, but must be deleted.")
	}
}

func checkFields(t testing.TB, actual, original any, fieldName string) {
	t.Helper()

	if actual != original {
		t.Errorf("%s original %v, actual %v", fieldName, original, actual)
	}
}

func checkJSONFields(t testing.TB, actual, original json.RawMessage, fieldName string) {
	t.Helper()

	if !jsonEqual(actual, original) {
		t.Errorf("%s original %s, actual %s", fieldName, original, actual)
	}
}

func jsonEqual(a, b json.RawMessage) bool {
	var left, right any

	if json.Unmarshal(a, &left) != nil {
		return false
	}

	if json.Unmarshal(b, &right) != nil {
		return false
	}

	return reflect.DeepEqual(left, right)
}
