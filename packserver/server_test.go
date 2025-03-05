package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// TestPackChange tests the packChange function with different inputs.
func TestPackChange(t *testing.T) {
	testCases := []struct {
		packs    []int
		order    int
		expected map[string]int
	}{
		{packs: []int{250, 500, 1000, 2000, 5000}, order: 1, expected: map[string]int{"250": 1}},
		{packs: []int{250, 500, 1000, 2000, 5000}, order: 250, expected: map[string]int{"250": 1}},
		{packs: []int{250, 500, 1000, 2000, 5000}, order: 251, expected: map[string]int{"500": 1}},
		{packs: []int{250, 500, 1000, 2000, 5000}, order: 501, expected: map[string]int{"250": 1, "500": 1}},
		{packs: []int{250, 500, 1000, 2000, 5000}, order: 12001, expected: map[string]int{"250": 1, "2000": 1, "5000": 2}},
	}

	for _, tc := range testCases {
		result := packChange(tc.packs, tc.order)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("packChange(%v, %d) = %v; expected %v", tc.packs, tc.order, result, tc.expected)
		}
	}
}

// TestPackChangeHandler tests the HTTP handler with valid input.
func TestPackChangeHandler(t *testing.T) {
	reqBody := Request{
		Packs: []int{250, 500, 1000, 2000, 5000},
		Order: 251,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Create a new POST request with the JSON body.
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	packChangeHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var resp map[string]int
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expected := map[string]int{"500": 1}
	if !reflect.DeepEqual(resp, expected) {
		t.Errorf("Expected response %v, got %v", expected, resp)
	}
}

// TestPackChangeHandlerInvalidOrder tests the handler with an invalid order (<= 0).
func TestPackChangeHandlerInvalidOrder(t *testing.T) {
	reqBody := Request{
		Packs: []int{250, 500, 1000, 2000, 5000},
		Order: 0, // invalid: order must be > 0
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	packChangeHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid order, got %d", rr.Code)
	}
}

// TestPackChangeHandlerInvalidPack tests the handler with an invalid pack size (<= 0).
func TestPackChangeHandlerInvalidPack(t *testing.T) {
	reqBody := Request{
		Packs: []int{250, -500, 1000}, // invalid: -500 is not allowed
		Order: 500,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	packChangeHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid pack size, got %d", rr.Code)
	}
}

// TestPackChangeHandlerNonNumeric tests the handler with non-numeric JSON fields.
func TestPackChangeHandlerNonNumeric(t *testing.T) {
	// Since the Request struct defines Order as int and Packs as []int,
	// passing a non-numeric value should trigger a JSON decode error.
	jsonStr := `{"packs": [250, 500, 1000], "order": "abc"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	packChangeHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for non-numeric order, got %d", rr.Code)
	}
}
