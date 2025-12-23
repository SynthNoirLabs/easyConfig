# Task: Dotfiles Sync (GitOps)

## 🎯 Objective
Provide a robust backup and sync mechanism for agent configurations using a private Git repository.

## 📝 Description
Developers work across multiple machines. They want their Claude settings and Gemini aliases to be the same everywhere. This task implements a simple "GitOps" flow.

## ✅ Requirements
1.  **Repo Initialization**:
    -   Button: "Initialize Backup Repo".
    -   Creates a hidden git repo in `~/.config/easyconfig/backup`.
2.  **Symlinking Strategy**:
    -   Move original config files to the backup folder.
    -   Create symlinks from the original locations to the backup folder.
3.  **Sync Logic**:
    -   "Push Changes": Commit and push to a user-provided remote URL.
    -   "Pull Changes": Pull and re-apply symlinks.

## 🛠️ Technical Implementation
-   **Backend**: `pkg/sync/gitops.go`.
-   **Frontend**: "Sync" status indicator in the footer.
