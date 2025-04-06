package common

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/google/uuid"
)

func RoundToNDecimals(value float64, n int) float64 {
	scale := math.Pow(10, float64(n))
	return math.Round(value*scale) / scale
}

func WriteJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// pretty-print JSON
	out, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		http.Error(w, `{"error":{"code":5000,"message":"Failed to serialize response"}}`, http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(out)
}

func IsValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

func GetIdempotencyKey(r *http.Request) string {
	return r.Header.Get("Idempotency-Key")
}
