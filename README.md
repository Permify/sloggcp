# sloggcp

A `log/slog` handler for Google Cloud Logging, enabling seamless integration with Google Cloud's logging infrastructure.

## Quick Start

### Installation

Ensure you have the Google Cloud SDK installed and authenticated. Then, add `sloggcp` and necessary dependencies:

```bash
go get github.com/Permify/sloggcp
```

### Usage
Below is a minimal example of configuring and using sloggcp to send logs to Google Cloud Logging.

```go
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Permify/sloggcp"
)

func main() {
	// Set up context and Google Cloud project ID
	ctx := context.Background()
	projectID := "YOUR_PROJECT_ID"
	logName := "my-application-log"

	// Create a new GoogleCloudSlogHandler
	handler, err := sloggcp.NewGoogleCloudSlogHandler(ctx, projectID, logName, slog.LevelInfo)
	if err != nil {
		slog.Error("Failed to create Google Cloud slog handler", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer handler.Close()

	// Set up the slog logger with the Google Cloud handler
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Log messages at different levels
	logger.Info("Application started", slog.String("env", "production"))
	logger.Debug("Debugging info", slog.String("module", "auth"))
	logger.Warn("Potential issue detected", slog.String("component", "database"))
	logger.Error("An error occurred", slog.String("error", "database connection failed"))

	// Optional: Close the handler when done to ensure logs are flushed
	if err := handler.Close(); err != nil {
		slog.Error("Failed to close the Google Cloud slog handler", slog.String("error", err.Error()))
	}
}
```
### Explanation

- Project ID and Log Name: Set projectID and logName to your Google Cloud Project ID and desired log name.
- slog Handler: GoogleCloudSlogHandler sends logs directly to Google Cloud, handling attributes and severity levels automatically.
- Logging Levels: Log messages can be sent with various levels (Info, Debug, Warn, Error), automatically converting to the corresponding Google Cloud severity level.

### Environment Setup

Make sure the environment running this code has the necessary permissions to write to Google Cloud Logging. You can achieve this by:

1. Setting up authentication with a service account key.
2. Configuring the GOOGLE_APPLICATION_CREDENTIALS environment variable to point to the key file.

```bash
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/your/service-account-key.json"
```

### Testing Logs
After running the example, go to Google Cloud Logging and check for entries under the specified log name (my-application-log).
```bash
go test -v
```
