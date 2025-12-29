package config

import (
	"os/exec"
	"strings"
	"sync"
)

// GetAllProviderStatuses collects detailed health reports for all registered providers.
func (s *DiscoveryService) GetAllProviderStatuses() []ProviderStatusReport {
	var reports []ProviderStatusReport
	var wg sync.WaitGroup
	reportsChan := make(chan ProviderStatusReport, len(s.providers))

	for _, p := range s.providers {
		wg.Add(1)
		go func(provider Provider) {
			defer wg.Done()
			reportsChan <- s.checkProviderHealth(provider)
		}(p)
	}

	wg.Wait()
	close(reportsChan)

	for report := range reportsChan {
		reports = append(reports, report)
	}

	return reports
}

func (s *DiscoveryService) checkProviderHealth(provider Provider) ProviderStatusReport {
	report := ProviderStatusReport{
		ProviderName: provider.Name(),
	}

	// For now, we'll base the health check on the existing Discover method.
	// A more advanced implementation would check for binaries, versions, etc.
	items, err := provider.Discover("") // Use an empty project path for a general check
	if err != nil {
		report.Message = "Error during discovery."
		return report
	}

	if len(items) > 0 {
		report.Configured = true
		report.Valid = true // Placeholder for now
		report.Message = "Configuration files found."
	} else {
		report.Message = "No configuration files found."
	}

	// This is a simplified check. A real implementation would need to know
	// the binary name for each provider.
	// For now, let's assume the binary name is a lowercase version of the provider name.
	binaryName := strings.ToLower(strings.Split(provider.Name(), " ")[0])
	path, err := exec.LookPath(binaryName)
	if err == nil && path != "" {
		report.Installed = true
		// Try to get the version
		// This is highly dependent on the tool's CLI
		out, err := exec.Command(path, "--version").Output()
		if err == nil {
			report.Version = strings.TrimSpace(string(out))
		}
	}

	return report
}
