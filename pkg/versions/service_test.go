package versions

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRepo creates a temporary git repository for testing.
func setupTestRepo(t *testing.T) (repoPath, filePath string) {
	t.Helper()

	// Create a temporary directory for the repo
	repoPath = t.TempDir()

	// Initialize a new git repository
	_, err := git.PlainInit(repoPath, false)
	require.NoError(t, err)

	// Create a dummy file and commit it
	filePath = filepath.Join(repoPath, "test.txt")
	err = os.WriteFile(filePath, []byte("initial commit"), 0644)
	require.NoError(t, err)

	repo, err := git.PlainOpen(repoPath)
	require.NoError(t, err)

	w, err := repo.Worktree()
	require.NoError(t, err)

	_, err = w.Add("test.txt")
	require.NoError(t, err)

	_, err = w.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test Author",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)

	return repoPath, filePath
}

func TestVersionControlService_GetFileHistory(t *testing.T) {
	_, filePath := setupTestRepo(t)
	service := NewVersionControlService(nil)

	history, err := service.GetFileHistory(filePath)
	assert.NoError(t, err)
	assert.Len(t, history, 1)
	assert.Equal(t, "Initial commit", history[0].Message)
}

func TestVersionControlService_GetFileContentAtCommit(t *testing.T) {
	_, filePath := setupTestRepo(t)
	service := NewVersionControlService(nil)

	history, err := service.GetFileHistory(filePath)
	require.NoError(t, err)
	require.Len(t, history, 1)

	content, err := service.GetFileContentAtCommit(filePath, history[0].Hash)
	assert.NoError(t, err)
	assert.Equal(t, "initial commit", content)
}
