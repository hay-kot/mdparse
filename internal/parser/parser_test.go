package parser

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestIsJSONObject(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty string", "", false},
		{"whitespace only", "   \n\t  ", false},
		{"valid object", `{"key": "value"}`, true},
		{"valid array", `["item1", "item2"]`, false},
		{"object with whitespace", `  {"key": "value"}  `, true},
		{"markdown", "# Hello World", false},
		{"yaml frontmatter", "---\ntitle: Test\n---", false},
		{"invalid json", `{"key": "value"`, false},
		{"number", "123", false},
		{"string", `"hello"`, false},
		{"boolean", "true", false},
		{"null", "null", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsJSONObject([]byte(tt.input)); got != tt.want {
				t.Errorf("IsJSONObject(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestMarkdownToJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		bodyKey string
		pretty  bool
		wantErr bool
		check   func(t *testing.T, output []byte)
	}{
		{
			name:    "markdown without frontmatter",
			input:   "# Hello World\n\nThis is content.",
			bodyKey: "$body",
			pretty:  false,
			wantErr: false,
			check: func(t *testing.T, output []byte) {
				var result map[string]any
				if err := json.Unmarshal(output, &result); err != nil {
					t.Errorf("output is not valid JSON: %v", err)
				}
				if body, ok := result["$body"].(string); !ok || body != "# Hello World\n\nThis is content." {
					t.Errorf("unexpected body content: %v", result["$body"])
				}
			},
		},
		{
			name: "markdown with yaml frontmatter",
			input: `---
title: Test Post
author: John Doe
tags:
  - test
  - example
---

# Hello World

This is content.`,
			bodyKey: "$body",
			pretty:  false,
			wantErr: false,
			check: func(t *testing.T, output []byte) {
				var result map[string]any
				if err := json.Unmarshal(output, &result); err != nil {
					t.Errorf("output is not valid JSON: %v", err)
				}
				if title, ok := result["title"].(string); !ok || title != "Test Post" {
					t.Errorf("unexpected title: %v", result["title"])
				}
				if author, ok := result["author"].(string); !ok || author != "John Doe" {
					t.Errorf("unexpected author: %v", result["author"])
				}
				if body, ok := result["$body"].(string); !ok || body != "# Hello World\n\nThis is content." {
					t.Errorf("unexpected body: %v", result["$body"])
				}
			},
		},
		{
			name:    "custom body key",
			input:   "# Test",
			bodyKey: "content",
			pretty:  false,
			wantErr: false,
			check: func(t *testing.T, output []byte) {
				var result map[string]any
				if err := json.Unmarshal(output, &result); err != nil {
					t.Errorf("output is not valid JSON: %v", err)
				}
				if _, ok := result["content"]; !ok {
					t.Errorf("content key not found in output")
				}
			},
		},
		{
			name:    "pretty JSON formatting",
			input:   "# Test",
			bodyKey: "$body",
			pretty:  true,
			wantErr: false,
			check: func(t *testing.T, output []byte) {
				if !bytes.Contains(output, []byte("\n")) {
					t.Errorf("expected pretty formatted JSON with newlines")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := MarkdownToJSON([]byte(tt.input), tt.bodyKey, tt.pretty)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarkdownToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.check != nil && !tt.wantErr {
				tt.check(t, output)
			}
		})
	}
}

func TestJSONToMarkdown(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		bodyKey    string
		wantErr    bool
		wantOutput string
	}{
		{
			name:       "json without frontmatter",
			input:      `{"$body": "# Hello World"}`,
			bodyKey:    "$body",
			wantErr:    false,
			wantOutput: "# Hello World",
		},
		{
			name:    "json with frontmatter",
			input:   `{"title": "Test Post", "author": "John Doe", "$body": "# Hello World\n\nContent here."}`,
			bodyKey: "$body",
			wantErr: false,
			wantOutput: `---
author: John Doe
title: Test Post
---

# Hello World

Content here.`,
		},
		{
			name:    "custom body key",
			input:   `{"title": "Test", "content": "# Hello"}`,
			bodyKey: "content",
			wantErr: false,
			wantOutput: `---
title: Test
---

# Hello`,
		},
		{
			name:    "missing body key",
			input:   `{"title": "Test"}`,
			bodyKey: "$body",
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			input:   `{"title": "Test"`,
			bodyKey: "$body",
			wantErr: true,
		},
		{
			name:    "body is not a string",
			input:   `{"$body": 123}`,
			bodyKey: "$body",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := JSONToMarkdown([]byte(tt.input), tt.bodyKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONToMarkdown() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(output) != tt.wantOutput {
				t.Errorf("JSONToMarkdown() output mismatch\nGot:\n%s\n\nWant:\n%s", string(output), tt.wantOutput)
			}
		})
	}
}
