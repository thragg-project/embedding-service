package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/fr33dman/go-template/internal/logger"
)

func TestLoggerWritesJSON(t *testing.T) {
	var out bytes.Buffer

	log, err := logger.NewWithWriter(&out, "info", "test-version")
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}

	log.InfoContext(context.Background(), "test message", "key", "value")

	var entry map[string]any
	if err := json.Unmarshal(out.Bytes(), &entry); err != nil {
		t.Fatalf("log output is not json: %v", err)
	}

	if entry["msg"] != "test message" {
		t.Fatalf("expected msg %q, got %v", "test message", entry["msg"])
	}
	if entry["key"] != "value" {
		t.Fatalf("expected key %q, got %v", "value", entry["key"])
	}
	if entry["version"] != "test-version" {
		t.Fatalf("expected version %q, got %v", "test-version", entry["version"])
	}
}

func TestLoggerRejectsUnknownLevel(t *testing.T) {
	var out bytes.Buffer

	_, err := logger.NewWithWriter(&out, "verbose", "test-version")
	if err == nil {
		t.Fatal("expected error")
	}
}
