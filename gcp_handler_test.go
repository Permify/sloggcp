package sloggcp

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"log/slog"

	"github.com/joho/godotenv"
)

func TestGoogleCloudSlogHandler(t *testing.T) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	// Set up context and Google Cloud project ID
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		t.Fatal("GOOGLE_CLOUD_PROJECT environment variable is not set")
	}

	// Initialize GoogleCloudSlogHandler
	logName := "test-log"
	handler, err := NewGoogleCloudSlogHandler(ctx, projectID, logName, slog.LevelInfo)
	if err != nil {
		t.Fatalf("failed to initialize GoogleCloudSlogHandler: %v", err)
	}
	defer func() {
		if err := handler.Close(); err != nil {
			t.Fatalf("failed to close handler: %v", err)
		}
	}()

	// Define test cases
	testCases := []struct {
		level   slog.Level
		message string
		attrs   []slog.Attr
	}{
		{slog.LevelDebug, "Debug message", []slog.Attr{slog.String("key1", "value1")}},
		{slog.LevelInfo, "Info message", []slog.Attr{slog.String("key2", "value2")}},
		{slog.LevelWarn, "Warning message", []slog.Attr{slog.Int("key3", 123)}},
		{slog.LevelError, "Error message", []slog.Attr{slog.Float64("key4", 3.14)}},
	}

	// Execute each test case
	for _, tc := range testCases {
		// Create a record for each test case
		record := slog.Record{
			Time:    time.Now(),
			Level:   tc.level,
			Message: tc.message,
		}

		for _, attr := range tc.attrs {
			record.AddAttrs(attr)
		}

		// Handle the record
		err := handler.Handle(ctx, record)
		if err != nil {
			t.Fatalf("unexpected error in Handle: %v", err)
		}
	}

	fmt.Printf("All logs written successfully to log: %s\n", logName)
}
