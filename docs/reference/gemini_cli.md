# Gemini CLI Configuration Reference

**Source:** https://geminicli.com/docs/cli/settings/

## Settings File Locations
*   **User settings:** `~/.gemini/settings.json`
*   **Workspace settings:** `./.gemini/settings.json` (Overrides user)

## Format: JSON

## Key Settings Categories

### 1. UI & Appearance
```json
"ui": {
  "hideWindowTitle": false,
  "showStatusInTitle": true,
  "hideTips": false,
  "useFullWidth": true
}
```

### 2. Context & Files
```json
"context": {
  "fileFiltering": {
    "respectGitIgnore": true,
    "respectGeminiIgnore": true,
    "enableRecursiveFileSearch": true
  },
  "discoveryMaxDirs": 200
}
```

### 3. Tools & Security
```json
"tools": {
  "useRipgrep": true,
  "autoAccept": false
},
"security": {
  "disableYoloMode": false,
  "blockGitExtensions": false
}
```

### 4. Model Parameters
```json
"model": {
  "maxSessionTurns": -1,
  "compressionThreshold": 0.2
}
```

## CLI Commands
*   `/settings`: Opens interactive settings dialog.
*   `/init`: Generates a `GEMINI.md` context file.

## Example `settings.json`
```json
{
  "ui": {
    "showLineNumbers": true
  },
  "context": {
    "fileFiltering": {
      "respectGitIgnore": true
    }
  },
  "general": {
    "disableAutoUpdate": false
  }
}
```