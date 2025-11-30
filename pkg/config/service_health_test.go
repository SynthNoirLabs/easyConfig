package config

import (
	"testing"
	"time"
)

// MockProvider allows us to simulate different provider states
type MockProvider struct {
	name    string
	healthy bool
	files   []Item
}

func (m *MockProvider) Name() string {
	return m.name
}

func (m *MockProvider) Discover(_ string) ([]Item, error) {
	return m.files, nil
}

func (m *MockProvider) Create(_ Scope, _ string) (string, error) {
	return "/mock/path", nil
}

func (m *MockProvider) CheckStatus() ProviderStatus {
	status := ProviderStatus{
		ProviderName: m.name,
		LastChecked:  time.Now().Format(time.RFC3339),
	}
	if m.healthy {
		status.Health = StatusHealthy
		status.StatusMessage = "Mock is healthy"
	} else {
		status.Health = StatusUnhealthy
		status.StatusMessage = "Mock is unhealthy"
	}
	status.DiscoveredFiles = m.files
	return status
}

func TestGetProviderStatuses(t *testing.T) {
	// Setup
	ds := NewDiscoveryService(nil) // logger is nil for test simplicity
	ds.providers = []Provider{
		&MockProvider{name: "HealthyProvider", healthy: true, files: []Item{{Name: "file1"}}},
		&MockProvider{name: "UnhealthyProvider", healthy: false},
		&MockProvider{name: "AnotherHealthy", healthy: true},
	}

	// Execute
	statuses := ds.GetProviderStatuses()

	// Verify
	if len(statuses) != 3 {
		t.Errorf("Expected 3 statuses, got %d", len(statuses))
	}

	// Check HealthyProvider
	if statuses[0].ProviderName != "HealthyProvider" || statuses[0].Health != StatusHealthy {
		t.Errorf("Expected HealthyProvider to be healthy, got %+v", statuses[0])
	}
	if len(statuses[0].DiscoveredFiles) != 1 {
		t.Errorf("Expected HealthyProvider to have 1 discovered file, got %d", len(statuses[0].DiscoveredFiles))
	}

	// Check UnhealthyProvider
	if statuses[1].ProviderName != "UnhealthyProvider" || statuses[1].Health != StatusUnhealthy {
		t.Errorf("Expected UnhealthyProvider to be unhealthy, got %+v", statuses[1])
	}

	// Check AnotherHealthy
	if statuses[2].ProviderName != "AnotherHealthy" || statuses[2].Health != StatusHealthy {
		t.Errorf("Expected AnotherHealthy to be healthy, got %+v", statuses[2])
	}
}
