package postgres_test

import (
	"buckwheat-queue/internal/test"
	"testing"
)

func TestNewDB(t *testing.T) {
	cfg := test.GetConfig(t)
	test.GetDB(t, cfg)
}
