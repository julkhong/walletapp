package common

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoundToNDecimals(t *testing.T) {
	tests := []struct {
		value    float64
		decimals int
		expected float64
	}{
		{123.456789, 2, 123.46},
		{123.456789, 3, 123.457},
		{123.451, 2, 123.45},
		{123.4, 0, 123},
		{0.00009, 4, 0.0001},
	}

	for _, tt := range tests {
		got := RoundToNDecimals(tt.value, tt.decimals)
		if got != tt.expected {
			t.Errorf("RoundToNDecimals(%f, %d) = %f; want %f", tt.value, tt.decimals, got, tt.expected)
		}
	}
}

func TestWriteJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	payload := map[string]string{"message": "hello world"}

	WriteJSON(rec, http.StatusOK, payload)

	res := rec.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Logf("Warning: failed to close response body: %v", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, res.StatusCode)
	}

	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected content-type 'application/json'; got %s", ct)
	}

	// Check response body
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(res.Body)
	got := buf.String()

	expected := `{
  "message": "hello world"
}`

	if got != expected {
		t.Errorf("unexpected JSON body:\n got: %s\nwant: %s", got, expected)
	}
}

func TestIsValidUUID(t *testing.T) {
	valid := "550e8400-e29b-41d4-a716-446655440000"
	invalid := "not-a-uuid"

	if !IsValidUUID(valid) {
		t.Errorf("expected valid UUID for %s", valid)
	}
	if IsValidUUID(invalid) {
		t.Errorf("expected invalid UUID for %s", invalid)
	}
}

func TestGetIdempotencyKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Idempotency-Key", "xyz-123")

	key := GetIdempotencyKey(req)
	if key != "xyz-123" {
		t.Errorf("expected idempotency key 'xyz-123'; got '%s'", key)
	}
}
