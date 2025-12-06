package versions

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// VersionControlService provides methods for interacting with git history.
type VersionControlService struct {
	logger *slog.Logger
}

// NewVersionControlService creates a new service.
func NewVersionControlService(logger *slog.Logger) *VersionControlService {
	if logger == nil {
		logger = slog.Default()
	}
	return &VersionControlService{logger: logger}
}

// GetFileHistory retrieves the commit history for a specific file.
func (s *VersionControlService) GetFileHistory(filePath string) ([]CommitInfo, error) {
	repo, err := s.findRepo(filePath)
	if err != nil {
		// If it's not a git repo, that's fine, just return no history.
		if err == git.ErrRepositoryNotExists {
			return []CommitInfo{}, nil
		}
		return nil, fmt.Errorf("could not find git repository for file %s: %w", filePath, err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("could not get worktree: %w", err)
	}
	relativePath, err := filepath.Rel(worktree.Filesystem.Root(), filePath)
	if err != nil {
		return nil, fmt.Errorf("could not determine relative path for %s: %w", filePath, err)
	}
	gitStylePath := filepath.ToSlash(relativePath)

	logOptions := &git.LogOptions{
		FileName: &gitStylePath,
		Order:    git.LogOrderCommitterTime,
	}

	cIter, err := repo.Log(logOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get git log for %s: %w", filePath, err)
	}

	var history []CommitInfo
	err = cIter.ForEach(func(c *object.Commit) error {
		history = append(history, CommitInfo{
			Hash:    c.Hash.String(),
			Author:  c.Author.Name,
			Message: c.Message,
			Date:    c.Author.When.Format(time.RFC3339),
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate through commits: %w", err)
	}

	return history, nil
}

// GetFileContentAtCommit retrieves the content of a file at a specific commit.
func (s *VersionControlService) GetFileContentAtCommit(filePath, commitHash string) (string, error) {
	repo, err := s.findRepo(filePath)
	if err != nil {
		return "", fmt.Errorf("could not find git repository for file %s: %w", filePath, err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("could not get worktree: %w", err)
	}
	relativePath, err := filepath.Rel(worktree.Filesystem.Root(), filePath)
	if err != nil {
		return "", fmt.Errorf("could not determine relative path for %s: %w", filePath, err)
	}
	gitStylePath := filepath.ToSlash(relativePath)

	hash := plumbing.NewHash(commitHash)
	commit, err := repo.CommitObject(hash)
	if err != nil {
		return "", fmt.Errorf("could not find commit %s: %w", commitHash, err)
	}

	tree, err := commit.Tree()
	if err != nil {
		return "", fmt.Errorf("could not get tree for commit %s: %w", commitHash, err)
	}

	file, err := tree.File(gitStylePath)
	if err != nil {
		return "", fmt.Errorf("could not find file %s in commit %s: %w", gitStylePath, commitHash, err)
	}

	reader, err := file.Reader()
	if err != nil {
		return "", fmt.Errorf("could not read file %s: %w", gitStylePath, err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read content of %s: %w", gitStylePath, err)
	}

	return string(content), nil
}

// findRepo finds the git repository containing the given path.
func (s *VersionControlService) findRepo(path string) (*git.Repository, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(absPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("path does not exist: %s", dir)
	}

	repo, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{
		DetectDotGit: true,
	})

	if err != nil {
		return nil, err
	}

	return repo, nil
}
