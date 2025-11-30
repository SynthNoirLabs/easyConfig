package marketplaces

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSmitheryClient_FetchPopularServers_Success(t *testing.T) {
	// Mock server response
	mockPackages := []MCPPackage{
		{Name: "test-pkg", Description: "Test Package", Version: "1.0.0", Author: "Tester", URL: "http://example.com"},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/packages/popular" {
			t.Errorf("Expected path /packages/popular, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(mockPackages); err != nil {
			t.Errorf("Failed to encode mock response: %v", err)
		}
	}))
	defer server.Close()

	client := NewSmitheryClient()
	client.BaseURL = server.URL
	client.HTTPClient = server.Client()

	packages, err := client.FetchPopularServers()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(packages) != 1 {
		t.Errorf("Expected 1 package, got %d", len(packages))
	}
	if packages[0].Name != "test-pkg" {
		t.Errorf("Expected package name 'test-pkg', got %s", packages[0].Name)
	}
}

func TestSmitheryClient_FetchPopularServers_Fallback(t *testing.T) {
	// Mock server error (404)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewSmitheryClient()
	client.BaseURL = server.URL
	client.HTTPClient = server.Client()

	packages, err := client.FetchPopularServers()
	if err != nil {
		t.Fatalf("Expected no error (fallback), got %v", err)
	}

	// Should return mock data (4 items)
	if len(packages) != 4 {
		t.Errorf("Expected 4 fallback packages, got %d", len(packages))
	}
}

func TestSmitheryClient_FetchPopularServers_NetworkError(t *testing.T) {
	// Client with no server (connection refused)
	client := NewSmitheryClient()
	client.BaseURL = "http://localhost:12345" // Unused port

	_, err := client.FetchPopularServers()
	if err == nil {
		t.Error("Expected network error, got nil")
	}
}
