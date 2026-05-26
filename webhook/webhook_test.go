package webhook

import (
	"strings"
	"testing"
	"time"
)

func TestSign(t *testing.T) {
	payload := []byte(`{"event":"test"}`)
	secret := "whsec_test123"
	ts := time.Unix(1700000000, 0)

	sig := Sign(payload, secret, ts)
	if !strings.HasPrefix(sig, "v1=") {
		t.Errorf("signature should start with v1=, got %q", sig)
	}

	sig2 := Sign(payload, secret, ts)
	if sig != sig2 {
		t.Error("same inputs should produce same signature")
	}

	sig3 := Sign(payload, "different-secret", ts)
	if sig == sig3 {
		t.Error("different secrets should produce different signatures")
	}
}

func TestVerify(t *testing.T) {
	payload := []byte(`{"event":"test"}`)
	secret := "whsec_test123"
	ts := time.Now()

	sig := Sign(payload, secret, ts)

	if !Verify(payload, secret, sig, ts, 5*time.Minute) {
		t.Error("valid signature should verify")
	}

	if Verify(payload, "wrong-secret", sig, ts, 5*time.Minute) {
		t.Error("wrong secret should fail verification")
	}

	if Verify(payload, secret, "v1=bad", ts, 5*time.Minute) {
		t.Error("wrong signature should fail verification")
	}
}

func TestVerifyTimestampTolerance(t *testing.T) {
	payload := []byte(`{"event":"test"}`)
	secret := "whsec_test123"
	oldTs := time.Now().Add(-10 * time.Minute)

	sig := Sign(payload, secret, oldTs)

	if Verify(payload, secret, sig, oldTs, 5*time.Minute) {
		t.Error("expired timestamp should fail with tolerance")
	}

	// time.Duration(0) is strict (no drift allowed), so a 10-minute-old
	// timestamp must NOT verify.
	if Verify(payload, secret, sig, oldTs, 0) {
		t.Error("zero tolerance should be strict and reject any drift")
	}

	// NoToleranceCheck disables the drift check entirely.
	if !Verify(payload, secret, sig, oldTs, NoToleranceCheck) {
		t.Error("NoToleranceCheck should skip timestamp check")
	}

	// Strict (0) on a fresh timestamp: still allow same-second.
	fresh := time.Now()
	freshSig := Sign(payload, secret, fresh)
	if !Verify(payload, secret, freshSig, fresh, NoToleranceCheck) {
		t.Error("fresh sig with NoToleranceCheck should verify")
	}
}

func TestVerifyMultiScheme(t *testing.T) {
	payload := []byte(`{"event":"test"}`)
	secret := "whsec_test123"
	ts := time.Now()

	v1 := Sign(payload, secret, ts) // "v1=..."

	// Multi-scheme header: v0 first, then v1.
	multi := "v0=legacysig," + v1
	if !Verify(payload, secret, multi, ts, 5*time.Minute) {
		t.Error("multi-scheme header should verify on v1 segment")
	}

	// v1 first, then unknown scheme.
	multi2 := v1 + ",v2=future"
	if !Verify(payload, secret, multi2, ts, 5*time.Minute) {
		t.Error("multi-scheme header with v1 first should verify")
	}

	// No v1 segment at all.
	if Verify(payload, secret, "v0=legacy,v2=future", ts, 5*time.Minute) {
		t.Error("header without v1 segment must not verify")
	}
}

func TestVerifyWithRotation(t *testing.T) {
	payload := []byte(`{"event":"test"}`)
	primary := "whsec_primary"
	secondary := "whsec_secondary"
	ts := time.Now()

	sigPrimary := Sign(payload, primary, ts)
	sigSecondary := Sign(payload, secondary, ts)

	if !VerifyWithRotation(payload, primary, secondary, sigPrimary, ts, 5*time.Minute) {
		t.Error("primary secret should verify")
	}

	if !VerifyWithRotation(payload, primary, secondary, sigSecondary, ts, 5*time.Minute) {
		t.Error("secondary secret should verify")
	}

	if VerifyWithRotation(payload, primary, secondary, "v1=bad", ts, 5*time.Minute) {
		t.Error("invalid signature should fail")
	}

	if VerifyWithRotation(payload, primary, "", sigSecondary, ts, 5*time.Minute) {
		t.Error("empty secondary should not verify secondary signature")
	}
}

func TestGenerateSecret(t *testing.T) {
	secret, err := GenerateSecret()
	if err != nil {
		t.Fatalf("GenerateSecret() error: %v", err)
	}

	if !strings.HasPrefix(secret, "whsec_") {
		t.Errorf("secret should start with whsec_, got %q", secret)
	}

	if len(secret) != 6+64 {
		t.Errorf("secret length = %d, want %d (6 prefix + 64 hex)", len(secret), 6+64)
	}

	secret2, _ := GenerateSecret()
	if secret == secret2 {
		t.Error("two generated secrets should not be equal")
	}
}
