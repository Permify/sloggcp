package sloggcp

import (
	"context"
	"fmt"
	"os"
	"testing"

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
	handler := NewGoogleCloudSlogHandler(ctx, projectID, logName, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	defer handler.Close()

	// Set the handler for slog
	slog.SetDefault(slog.New(handler))

	// Example log entries
	slog.Info("Starting application", "version", "1.0")
	slog.Debug("Debug", "debug", "sample debug")
	slog.Warn("This is a warning message", "component", "main")
	slog.Error("An error occurred", "error", "sample error")

	fmt.Printf("All logs written successfully to log: %s\n", logName)
}
