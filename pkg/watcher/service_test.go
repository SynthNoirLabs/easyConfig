package watcher

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWatcherService(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("initial"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Initialize service
	service := NewService()
	if service == nil {
		t.Fatal("Failed to create service")
	}
	service.SetEmitter(func(context.Context, string, ...interface{}) {})
	defer service.Close()

	// Test Add
	if err := service.Add(testFile); err != nil {
		t.Fatalf("Failed to add file to watcher: %v", err)
	}

	// Verify it's tracked
	service.mu.Lock()
	if !service.watched[testFile] {
		t.Error("File should be marked as watched")
	}
	service.mu.Unlock()

	// Test Remove
	if err := service.Remove(testFile); err != nil {
		t.Fatalf("Failed to remove file from watcher: %v", err)
	}

	// Verify it's untracked
	service.mu.Lock()
	if service.watched[testFile] {
		t.Error("File should be marked as unwatched")
	}
	service.mu.Unlock()
}

func TestWatcherService_AddNonExistent(t *testing.T) {
	service := NewService()
	if service == nil {
		t.Fatal("Failed to create service")
	}
	service.SetEmitter(func(context.Context, string, ...interface{}) {})
	defer service.Close()

	// Try to watch a non-existent file
	err := service.Add("/path/to/non/existent/file")
	if err == nil {
		t.Error("Expected error when adding non-existent file, got nil")
	}
}

func TestWatcherService_StartAndClose(t *testing.T) {
	service := NewService()
	if service == nil {
		t.Fatal("Failed to create service")
	}
	service.SetEmitter(func(context.Context, string, ...interface{}) {})

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Start the service
	service.Start(ctx)

	// Give it a moment to start
	time.Sleep(50 * time.Millisecond)

	// Add a file to watch
	tempDir := t.TempDir()
	// defer os.RemoveAll(tempDir) // t.TempDir automatically cleans up

	testFile := filepath.Join(tempDir, "watch_test.txt")
	if err := os.WriteFile(testFile, []byte("initial"), 0644); err != nil {
		t.Fatalf("Failed to write initial file: %v", err)
	}
	if err := service.Add(testFile); err != nil {
		t.Fatalf("Failed to add file to watcher: %v", err)
	}

	// Modify the file to trigger event
	time.Sleep(50 * time.Millisecond)
	if err := os.WriteFile(testFile, []byte("modified"), 0644); err != nil {
		t.Fatalf("Failed to write modified file: %v", err)
	}

	// Give it a moment to process event
	time.Sleep(100 * time.Millisecond)

	// Cancel context to stop loop
	cancel()

	// Give it a moment to stop
	time.Sleep(50 * time.Millisecond)

	// Close service
	service.Close()
}
