package config

import (
	"testing"
)

type MockProvider struct {
	name  string
	items []ConfigItem
}

func (m *MockProvider) Name() string {
	return m.name
}

func (m *MockProvider) Discover(projectPath string) ([]ConfigItem, error) {
	return m.items, nil
}

func TestDiscoveryService_DiscoverAll(t *testing.T) {
	ds := NewDiscoveryService()

	mock1 := &MockProvider{
		name: "Mock1",
		items: []ConfigItem{
			{Provider: "Mock1", Name: "Config1", Path: "/tmp/1", Scope: ScopeGlobal},
		},
	}
	ds.RegisterProvider(mock1)

	results, err := ds.DiscoverAll("/tmp")
	if err != nil {
		t.Fatalf("DiscoverAll failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 config, got %d", len(results))
	}
	if results[0].Provider != "Mock1" {
		t.Errorf("Expected provider Mock1, got %s", results[0].Provider)
	}
}
