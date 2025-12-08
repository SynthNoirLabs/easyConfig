package config

import (
	"os"
	"path/filepath"
	"testing"
)

type dummyProvider struct {
	path string
}

func (d *dummyProvider) Name() string { return "Dummy" }

func (d *dummyProvider) Discover(projectPath string) ([]Item, error) {
	return []Item{{
		Provider: d.Name(),
		Name:     "dummy",
		FileName: filepath.Base(d.path),
		Path:     d.path,
		Scope:    ScopeProject,
		Format:   "json",
		Exists:   true,
	}}, nil
}

func (d *dummyProvider) Create(scope Scope, projectPath string) (string, error) {
	return "", nil
}

func (d *dummyProvider) CheckStatus() ProviderStatus {
	return ProviderStatus{
		ProviderName: d.Name(),
		Health:       StatusUnknown,
	}
}

func (d *dummyProvider) GetWizard() Wizard {
	return nil
}

func TestSaveAndApplyProfile(t *testing.T) {
	tmp := t.TempDir()
	// isolate config dir
	if err := os.Setenv("XDG_CONFIG_HOME", tmp); err != nil {
		t.Fatalf("set XDG_CONFIG_HOME: %v", err)
	}
	if err := os.Setenv("HOME", tmp); err != nil {
		t.Fatalf("set HOME: %v", err)
	}
	defer func() {
		_ = os.Unsetenv("XDG_CONFIG_HOME")
	}()

	proj := filepath.Join(tmp, "proj")
	if err := os.MkdirAll(proj, 0o755); err != nil {
		t.Fatal(err)
	}

	cfgPath := filepath.Join(proj, "dummy.json")
	if err := os.WriteFile(cfgPath, []byte(`{"hello":"world"}`), 0o600); err != nil {
		t.Fatal(err)
	}

	svc := NewDiscoveryService(nil)
	svc.RegisterProvider(&dummyProvider{path: cfgPath})

	if err := svc.SaveProfile("test-profile", proj); err != nil {
		t.Fatalf("save profile: %v", err)
	}

	list, err := svc.ListProfiles()
	if err != nil {
		t.Fatalf("list profiles: %v", err)
	}
	if len(list) != 1 || list[0].Name != "test-profile" {
		t.Fatalf("unexpected profiles: %+v", list)
	}

	// mutate file then apply to restore
	if err := os.WriteFile(cfgPath, []byte(`{"hello":"mutated"}`), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := svc.ApplyProfile("test-profile"); err != nil {
		t.Fatalf("apply profile: %v", err)
	}

	out, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != `{"hello":"world"}` {
		t.Fatalf("expected restored content, got %s", string(out))
	}

	if err := svc.DeleteProfile("test-profile"); err != nil {
		t.Fatalf("delete profile: %v", err)
	}

	// ensure timestamps were set
	if list[0].UpdatedAt.IsZero() {
		t.Fatalf("updatedAt not set")
	}
}
