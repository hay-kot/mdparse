package commands

import (
	"path/filepath"
	"strings"
)

// isLikelyPath returns true if the string looks like a file path
func isLikelyPath(s string) bool {
	if strings.HasPrefix(s, "{") {
		return false
	}

	// Empty string is not a path
	if s == "" {
		return false
	}

	// Check for common path indicators
	if strings.HasPrefix(s, "/") || // Absolute path (Unix)
		strings.HasPrefix(s, "~") || // Home directory
		strings.HasPrefix(s, "./") || // Relative path
		strings.HasPrefix(s, "../") || // Parent directory
		strings.Contains(s, string(filepath.Separator)) || // Contains path separator
		filepath.IsAbs(s) { // Absolute path (cross-platform, includes Windows C:\)
		return true
	}

	// Check for file extensions (common markdown extensions)
	ext := filepath.Ext(s)
	if ext == ".md" || ext == ".markdown" || ext == ".txt" {
		return true
	}

	return false
}
