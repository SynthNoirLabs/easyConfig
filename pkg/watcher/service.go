package watcher

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Service handles file watching and emits events when files change
type Service struct {
	ctx     context.Context
	watcher *fsnotify.Watcher
	watched map[string]bool
	mu      sync.Mutex
}

// NewService creates a new watcher service
func NewService() *Service {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Failed to create watcher: %v", err)
		return nil
	}

	return &Service{
		watcher: w,
		watched: make(map[string]bool),
	}
}

// Start begins the watcher loop
func (s *Service) Start(ctx context.Context) {
	s.ctx = ctx
	go s.watchLoop()
}

// Close stops the watcher
func (s *Service) Close() {
	if s.watcher != nil {
		s.watcher.Close()
	}
}

// Add adds a file to the watcher
func (s *Service) Add(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.watcher == nil {
		return fmt.Errorf("watcher not initialized")
	}

	// fsnotify watches directories, so we watch the parent dir
	// In a real app, we'd filter events for just this file in the loop
	// but for simplicity here we will watch the file directly if supported
	// or the dir. fsnotify supports watching files directly on many OSs.
	// Let's try watching the file directly first.

	if s.watched[path] {
		return nil
	}

	if err := s.watcher.Add(path); err != nil {
		return err
	}

	s.watched[path] = true
	log.Printf("Started watching: %s", path)
	return nil
}

// Remove stops watching a file
func (s *Service) Remove(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.watcher == nil {
		return nil
	}

	if !s.watched[path] {
		return nil
	}

	if err := s.watcher.Remove(path); err != nil {
		// Ignore "can't remove non-existent watch" errors
		log.Printf("Error removing watch: %v", err)
	}

	delete(s.watched, path)
	return nil
}

func (s *Service) watchLoop() {
	if s.watcher == nil {
		return
	}

	for {
		select {
		case <-s.ctx.Done():
			return
		case event, ok := <-s.watcher.Events:
			if !ok {
				return
			}
			// We only care about Write events for now (content changed)
			if event.Has(fsnotify.Write) {
				log.Printf("File modified: %s", event.Name)
				// Emit event to frontend
				// We send the absolute path so frontend can match it
				absPath, err := filepath.Abs(event.Name)
				if err == nil {
					runtime.EventsEmit(s.ctx, "config:changed", absPath)
				}
			}
		case err, ok := <-s.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}
