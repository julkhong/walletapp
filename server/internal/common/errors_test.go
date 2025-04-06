package common

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, http.StatusBadRequest, 1001, "Invalid input")

	res := rec.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Logf("Warning: failed to close response body: %v", err)
		}
	}()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, res.StatusCode)
	}

	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected content-type 'application/json'; got %s", ct)
	}

	expectedJSON := `{"error":{"code":1001,"message":"Invalid input"}}`

	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(res.Body)
	got := buf.String()

	if got != expectedJSON+"\n" && got != expectedJSON {
		t.Errorf("unexpected response body:\n got: %s\nwant: %s", got, expectedJSON)
	}
}
