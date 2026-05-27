# Hivehook Go SDK

Official Go client for [Hivehook](https://hivehook.com), webhook infrastructure for modern teams (inbound and outbound).

Latest release: **0.1.1** on [pkg.go.dev](https://pkg.go.dev/github.com/hivehook/sdk-go).

## Install

```bash
go get github.com/hivehook/sdk-go
```

## Quick start

```go
package main

import (
	"context"
	"log"

	hivehook "github.com/hivehook/sdk-go"
)

func main() {
	client := hivehook.New(
		hivehook.WithBaseURL("http://localhost:8080"),
		hivehook.WithAPIKey("your-api-key"),
	)

	source, err := client.Sources.Create(context.Background(), &hivehook.CreateSourceInput{
		Name:         "Stripe production",
		Slug:         "stripe-prod",
		ProviderType: "stripe",
		VerifyConfig: map[string]any{"secret": "whsec_..."},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("created source %s. POST webhooks to /ingest/%s", source.ID, source.Slug)
}
```

## Webhook signature verification

The `webhook` sub-package verifies inbound HMAC signatures from Hivehook's outbound deliveries:

```go
import (
	"strconv"
	"time"

	"github.com/hivehook/sdk-go/webhook"
)

signature := req.Header.Get(webhook.HeaderSignature)
unix, err := strconv.ParseInt(req.Header.Get(webhook.HeaderTimestamp), 10, 64)
if err != nil {
    http.Error(w, "invalid timestamp header", http.StatusBadRequest)
    return
}
ok := webhook.Verify(body, "your-signing-secret", signature, time.Unix(unix, 0), 5*time.Minute)
```

## Documentation

See the full reference at [hivehook.com/docs](https://hivehook.com/docs).

## License

MIT. See [LICENSE](LICENSE).
