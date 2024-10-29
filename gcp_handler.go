/*
Copyright Permify Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sloggcp

import (
	"context"
	"fmt"

	"log/slog"

	"cloud.google.com/go/logging"
)

// GoogleCloudSlogHandler wraps Google Cloud Logging's Logger for use with slog.
type GoogleCloudSlogHandler struct {
	client *logging.Client
	logger *logging.Logger
	level  slog.Level
}

// NewGoogleCloudSlogHandler initializes a new GoogleCloudSlogHandler.
func NewGoogleCloudSlogHandler(ctx context.Context, projectID, logName string, level slog.Level) (*GoogleCloudSlogHandler, error) {
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create logging client: %w", err)
	}
	return &GoogleCloudSlogHandler{
		client: client,
		logger: client.Logger(logName),
		level:  level,
	}, nil
}

// Handle adapts slog.Record entries to Google Cloud Logging entries.
func (h *GoogleCloudSlogHandler) Handle(ctx context.Context, r slog.Record) error {
	// Convert slog.Level to Google Cloud Logging severity
	var severity logging.Severity
	switch r.Level {
	case slog.LevelDebug:
		severity = logging.Debug
	case slog.LevelInfo:
		severity = logging.Info
	case slog.LevelWarn:
		severity = logging.Warning
	case slog.LevelError:
		severity = logging.Error
	default:
		severity = logging.Default
	}

	// Construct the payload with message and additional attributes
	payload := map[string]interface{}{
		"message": r.Message,
		"level":   r.Level.String(),
		"time":    r.Time,
	}

	// Add attributes from slog fields
	r.Attrs(func(a slog.Attr) bool {
		payload[a.Key] = a.Value
		return true
	})

	// Send log entry to Google Cloud Logging
	h.logger.Log(logging.Entry{
		Payload:  payload,
		Severity: severity,
	})

	return nil
}

// Close closes the Google Cloud Logging client.
func (h *GoogleCloudSlogHandler) Close() error {
	return h.client.Close()
}
