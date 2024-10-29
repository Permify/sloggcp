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
	logger      *logging.Logger
	client      *logging.Client
	opts        *slog.HandlerOptions
	groupPrefix string
	attrs       []slog.Attr
}

var _ slog.Handler = &GoogleCloudSlogHandler{}

// NewGoogleCloudSlogHandler initializes a new GoogleCloudSlogHandler.
func NewGoogleCloudSlogHandler(ctx context.Context, projectID, logName string, opts *slog.HandlerOptions) *GoogleCloudSlogHandler {
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		return nil
	}
	return &GoogleCloudSlogHandler{
		client: client,
		logger: client.Logger(logName),
		opts:   opts,
	}
}

func (h *GoogleCloudSlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

// Handle adapts slog.Record entries to Google Cloud Logging entries.
func (h *GoogleCloudSlogHandler) Handle(ctx context.Context, r slog.Record) error {
	// Check if the log level is enabled
	if !h.Enabled(ctx, r.Level) {
		return nil
	}

	// Convert slog.Level to Google Cloud Logging severity
	severity := h.mapSeverity(r.Level)

	// Construct the payload with message, time, and additional attributes
	payload := map[string]interface{}{
		"message": r.Message,
		"time":    r.Time,
	}

	// Add attributes from slog fields
	r.Attrs(func(a slog.Attr) bool {
		payload[a.Key] = h.formatAttrValue(a.Value)
		return true
	})

	// Send log entry to Google Cloud Logging and return any errors
	h.logger.Log(logging.Entry{
		Payload:  payload,
		Severity: severity,
	})

	return nil
}

func (h *GoogleCloudSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	for i, attr := range attrs {
		attrs[i] = withGroupPrefix(h.groupPrefix, attr)
	}

	return &GoogleCloudSlogHandler{
		client:      h.client,
		logger:      h.logger,
		opts:        h.opts,
		groupPrefix: h.groupPrefix,
		attrs:       append(h.attrs, attrs...),
	}
}

func (h *GoogleCloudSlogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	prefix := name + "."
	if h.groupPrefix != "" {
		prefix = h.groupPrefix + prefix
	}

	return &GoogleCloudSlogHandler{
		client:      h.client,
		logger:      h.logger,
		opts:        h.opts,
		attrs:       h.attrs,
		groupPrefix: prefix,
	}
}

// mapSeverity converts slog.Level to Google Cloud Logging's Severity.
func (h *GoogleCloudSlogHandler) mapSeverity(level slog.Level) logging.Severity {
	switch level {
	case slog.LevelDebug:
		return logging.Debug
	case slog.LevelInfo:
		return logging.Info
	case slog.LevelWarn:
		return logging.Warning
	case slog.LevelError:
		return logging.Error
	default:
		return logging.Default
	}
}

// formatAttrValue formats attribute values for Google Cloud Logging.
func (h *GoogleCloudSlogHandler) formatAttrValue(value interface{}) interface{} {
	switch v := value.(type) {
	case string, int, int64, float64, bool:
		return v
	case error:
		return v.Error()
	default:
		return fmt.Sprintf("%v", v) // Fallback for unsupported types
	}
}

func withGroupPrefix(groupPrefix string, attr slog.Attr) slog.Attr {
	if groupPrefix != "" {
		attr.Key = groupPrefix + attr.Key
	}
	return attr
}

// Close closes the Google Cloud Logging client.
func (h *GoogleCloudSlogHandler) Close() error {
	return h.client.Close()
}
