package testutils

import (
	"encoding/json"
	"testing"
)

func MarshalJSON(t *testing.T, j map[string]interface{}) []byte {
	b, err := json.Marshal(j)
	if err != nil {
		t.Fatalf("failed to marshal json %+v", j)
	}
	return b
}

func UnmarshalJSON(t *testing.T, b []byte) map[string]interface{} {
	var j map[string]interface{}
	err := json.Unmarshal(b, &j)
	if err != nil {
		t.Fatalf("failed to unmarshal json %+v", b)
	}
	return j
}
