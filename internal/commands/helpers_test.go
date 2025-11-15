package commands

import "testing"

func TestIsLikelyPath(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty string", "", false},
		{"absolute unix path", "/path/to/file", true},
		{"home directory", "~/file.md", true},
		{"relative current", "./file.md", true},
		{"relative parent", "../file.md", true},
		{"markdown extension", "file.md", true},
		{"markdown extension alt", "file.markdown", true},
		{"txt extension", "file.txt", true},
		{"windows absolute", "C:\\path\\file.md", true},
		{"path with separator", "path/to/file", true},
		{"plain text no extension", "hello world", false},
		{"markdown content", "# Title", false},
		{"short word", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLikelyPath(tt.input); got != tt.want {
				t.Errorf("isLikelyPath(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
