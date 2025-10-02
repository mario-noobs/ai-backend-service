package helper

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAPI(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		
		// Verify authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got %s", auth)
		}
		
		// Verify query parameters
		if r.URL.Query().Get("param1") != "value1" {
			t.Errorf("Expected param1=value1, got %s", r.URL.Query().Get("param1"))
		}
		
		// Return test response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	// Test parameters
	params := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	jwt := "Bearer test-token"

	// Call the function
	response, err := GetAPI(server.URL, params, jwt)
	
	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if result["success"] != true {
		t.Errorf("Expected success=true, got %v", result["success"])
	}
}

func TestGetAPI_InvalidURL(t *testing.T) {
	params := map[string]string{}
	jwt := "Bearer test-token"
	
	// Test with invalid URL
	_, err := GetAPI("://invalid-url", params, jwt)
	
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}

func TestPostAPI(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		
		// Verify content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got %s", contentType)
		}
		
		// Verify authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got %s", auth)
		}
		
		// Return test response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	// Test data
	jsonData := map[string]interface{}{
		"name":  "test",
		"value": 123,
	}
	jwt := "Bearer test-token"

	// Call the function
	response, err := PostAPI(server.URL, jsonData, jwt)
	
	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if result["message"] != "success" {
		t.Errorf("Expected message='success', got %v", result["message"])
	}
}

func TestPostAPI_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	// Test with invalid JSON data (function type cannot be marshaled)
	jsonData := map[string]interface{}{
		"invalid": func() {},
	}
	jwt := "Bearer test-token"

	// Call the function
	_, err := PostAPI(server.URL, jsonData, jwt)
	
	if err == nil {
		t.Error("Expected error for invalid JSON data, got nil")
	}
}

func TestPostAPI_InvalidURL(t *testing.T) {
	jsonData := map[string]interface{}{
		"test": "data",
	}
	jwt := "Bearer test-token"
	
	// Test with invalid URL
	_, err := PostAPI("://invalid-url", jsonData, jwt)
	
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}
