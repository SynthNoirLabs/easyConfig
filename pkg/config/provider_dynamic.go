package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"easyConfig/pkg/util/paths"
	"gopkg.in/yaml.v3"
)

// --- Dynamic Provider ---

// DynamicProviderDef defines the structure of a provider.yaml file
type DynamicProviderDef struct {
	Name        string   `yaml:"name"`
	BinaryName  string   `yaml:"binaryName"`
	VersionArgs []string `yaml:"versionArgs"`
	Files       []struct {
		Name     string `yaml:"name"`
		FileName string `yaml:"fileName"`
		Scope    Scope  `yaml:"scope"`
		Format   Format `yaml:"format"`
	} `yaml:"files"`
}

// DynamicProvider is a generic provider configured from a file
type DynamicProvider struct {
	def *DynamicProviderDef
}

// NewDynamicProvider creates a new dynamic provider from a definition file
func NewDynamicProvider(defPath string) (*DynamicProvider, error) {
	data, err := os.ReadFile(defPath)
	if err != nil {
		return nil, err
	}

	var def DynamicProviderDef
	if err := yaml.Unmarshal(data, &def); err != nil {
		return nil, err
	}

	return &DynamicProvider{def: &def}, nil
}

func (p *DynamicProvider) Name() string {
	return p.def.Name
}

func (p *DynamicProvider) BinaryName() string {
	return p.def.BinaryName
}

func (p *DynamicProvider) VersionArgs() []string {
	return p.def.VersionArgs
}

func (p *DynamicProvider) Create(scope Scope, projectPath string) (string, error) {
	return "", fmt.Errorf("creating configs for dynamic providers is not yet supported")
}

func (p *DynamicProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := paths.GetHomeDir()
	configDir := paths.GetConfigDir(p.Name())

	for _, f := range p.def.Files {
		var path string
		switch f.Scope {
		case ScopeGlobal:
			if home == "" {
				continue
			}
			path = filepath.Join(home, f.FileName)
		case ScopeProject:
			if projectPath == "" {
				continue
			}
			path = filepath.Join(projectPath, f.FileName)
		case ScopeSystem:
			if configDir == "" {
				continue
			}
			path = filepath.Join(configDir, f.FileName)
		default:
			continue
		}

		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     f.Name,
				FileName: f.FileName,
				Path:     path,
				Scope:    f.Scope,
				Format:   f.Format,
				Exists:   true,
			})
		}
	}

	return items, nil
}

func (p *DynamicProvider) CheckStatus() ProviderStatus {
	status := ProviderStatus{
		ProviderName: p.Name(),
		LastChecked:  time.Now().Format(time.RFC3339),
	}

	files, _ := p.Discover("")
	status.DiscoveredFiles = files

	if len(files) > 0 {
		status.Health = StatusHealthy
		status.StatusMessage = "Configuration files found."
	} else {
		status.Health = StatusUnknown
		status.StatusMessage = "No configuration files found."
	}

	return status
}
