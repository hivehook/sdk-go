// Package webhook verifies inbound HMAC signatures from Hivehook's outbound
// webhook deliveries and signs payloads for tests.
package webhook

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// HeaderSignature is the request header carrying the v1=... HMAC signature.
const HeaderSignature = "X-Hivehook-Signature"

// HeaderTimestamp is the request header carrying the Unix-seconds timestamp.
const HeaderTimestamp = "X-Hivehook-Timestamp"

// HeaderMessageID is the request header carrying the unique delivery ID.
const HeaderMessageID = "X-Hivehook-Message-ID"

// NoToleranceCheck disables timestamp drift checking in Verify and
// VerifyWithRotation. Pass time.Duration(0) for strict (no drift allowed).
const NoToleranceCheck = time.Duration(-1)

// Sign produces a v1=<hex> signature header value over (timestamp + "." + payload).
func Sign(payload []byte, secret string, timestamp time.Time) string {
	ts := strconv.FormatInt(timestamp.Unix(), 10)
	msg := ts + "." + string(payload)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(msg))
	sig := hex.EncodeToString(mac.Sum(nil))

	return fmt.Sprintf("v1=%s", sig)
}

// Verify validates a signature header against payload using secret.
//
// The signature may be a single "v1=..." value or a multi-scheme header of the
// form "v0=...,v1=..."; the v1 segment is used.
//
// tolerance controls timestamp drift:
//   - time.Duration(0)      strict, timestamp must equal "now"
//   - NoToleranceCheck (-1) drift check disabled
//   - positive value        max absolute drift permitted
func Verify(payload []byte, secret string, signature string, timestamp time.Time, tolerance time.Duration) bool {
	if tolerance != NoToleranceCheck {
		age := time.Since(timestamp)
		if age < 0 {
			age = -age
		}
		if age > tolerance {
			return false
		}
	}

	v1, ok := extractV1(signature)
	if !ok {
		return false
	}

	expected := Sign(payload, secret, timestamp)
	expV1, _ := extractV1(expected)
	return hmac.Equal([]byte(expV1), []byte(v1))
}

// VerifyWithRotation validates against a primary secret and, on mismatch, a
// secondary secret (during key rotation). An empty secondary disables fallback.
func VerifyWithRotation(payload []byte, primary, secondary string, signature string, timestamp time.Time, tolerance time.Duration) bool {
	if Verify(payload, primary, signature, timestamp, tolerance) {
		return true
	}
	if secondary != "" {
		return Verify(payload, secondary, signature, timestamp, tolerance)
	}
	return false
}

// GenerateSecret returns a new "whsec_<32 random bytes hex>" signing secret.
func GenerateSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generating signing secret: %w", err)
	}
	return "whsec_" + hex.EncodeToString(b), nil
}

// extractV1 returns the hex bytes of the v1= segment from a possibly
// multi-scheme signature header (e.g. "v0=foo,v1=bar"). It returns ("", false)
// when no v1 segment is present.
func extractV1(sig string) (string, bool) {
	for _, part := range strings.Split(sig, ",") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "v1=") {
			return strings.TrimPrefix(part, "v1="), true
		}
	}
	return "", false
}
