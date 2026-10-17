package marketplace

import (
	"encoding/json"
	"testing"
)

func TestAPIErrorDecoding(t *testing.T) {
	var e APIError
	if err := json.Unmarshal([]byte(`{"code":"INVALID_ARGUMENT","message":"unknown realtime vendor"}`), &e); err != nil {
		t.Fatal(err)
	}
	if e.Code != "INVALID_ARGUMENT" || e.Message != "unknown realtime vendor" {
		t.Fatalf("unexpected API error: %+v", e)
	}
}
